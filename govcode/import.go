package main

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	c "govcode.org/common"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

func getStr(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func findOrCreateUser(user *c.User) int64 {

	var new_user c.User
	c.DB.Where("gh_id = ?", user.GhId).First(&new_user)

	if new_user.Id == 0 {
		new_user.Login = user.Login
		new_user.AvatarUrl = user.AvatarUrl
		new_user.GhId = user.GhId
		new_user.CommitCount = 0
		c.DB.Save(&new_user)
	}

	return new_user.Id
}

func runImport() (err error) {
	gh_key := os.Getenv("GH_KEY")

	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: gh_key},
	}

	client := github.NewClient(t.Client())

	orgs := getAccounts()
	importOrgs(orgs, client)

	return nil
}

func getAccounts() (orgs []string) {
	url := "http://registry.usa.gov/accounts.json?service_id=github"
	res, err := http.Get(url)
	c.PanicOn(err)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	c.PanicOn(err)

	var data AccountData

	err = json.Unmarshal(body, &data)
	c.PanicOn(err)

	orgs = append(orgs, []string{"MeasureAuthoringTool", "blue-button", "ngageoint", "afrl",
		"virtual-world-framework", "usnistgov", "govxteam", "commercegov",
		"ncats", "missioncommand"}...)

	for _, e := range data.Accounts {
		orgs = append(orgs, e.Account)
	}

	return orgs
}

func importOrgs(orgs []string, client *github.Client) {
	// Query github for org info and save to db if it doesn't exist
	var wg sync.WaitGroup
	for _, e := range orgs {
		wg.Add(1)

		go func(org_name string, wg *sync.WaitGroup) {
			gh_org, res, err := client.Organizations.Get(org_name)
			fmt.Println("Getting Org", org_name)
			fmt.Println(res)
			// Do not panic, probably 404 error
			if err != nil {
				fmt.Println(err)
				wg.Done()
				return
			}

			var org c.Organization

			c.DB.Where("login = ?", gh_org.Login).First(&org)

			if org.Id == 0 {
				org.Login = *gh_org.Login
				if gh_org.Name != nil {
					org.Name = *gh_org.Name
				}
				org.Ignore = false

				c.DB.Save(&org)
			}
			// Import repos
			importRepos(&org, client, 1)

			wg.Done()
		}(e, &wg)
	}
	wg.Wait()
}

func importRepos(org *c.Organization, client *github.Client, page int) {
	// Load the repos for a given org
	opt := &github.RepositoryListByOrgOptions{}
	opt.ListOptions = github.ListOptions{Page: page, PerPage: 100}

	repos, response, err := client.Repositories.ListByOrg(org.Login, opt)
	fmt.Println("Getting Repos", org.Login, page)
	fmt.Println(response)
	if err != nil {
		fmt.Println("Error with org: ", org.Login)
		return
	}

	// There are more pages, lets fetch the next one
	if response.NextPage > 0 {
		importRepos(org, client, response.NextPage)
	}

	concurrency := 4

	sem := make(chan bool, concurrency)

	for _, r := range repos {
		if !*r.Fork {
			var repo c.Repository
			c.DB.Where("name = ? and organization_id = ?", *r.Name, org.Id).First(&repo)

			if repo.Id == 0 {
				repo.GhId = int64(*r.ID)
				repo.Name = *r.Name
				if r.Description != nil {
					repo.Description = *r.Description
				}
				repo.OrganizationId = org.Id
			}
			repo.Forks = int64(*r.ForksCount)
			repo.Watchers = int64(*r.WatchersCount)
			repo.Stargazers = int64(*r.StargazersCount)
			repo.Size = int64(*r.Size)
			repo.OpenIssues = int64(*r.OpenIssuesCount)
			repo.Language = getStr(r.Language)

			c.DB.Save(&repo)

			sem <- true
			go importStats(&repo, org, client, &sem)

			importPulls(&repo, org, client, 1)
		}
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
}

func importStats(repo *c.Repository, org *c.Organization, client *github.Client, sem *chan bool) {
	stats, res, err := client.Repositories.ListContributorsStats(org.Login, repo.Name)
	fmt.Println("Getting stats", repo.Name)
	fmt.Println(res)
	if err != nil {
		if res != nil && res.StatusCode == 202 {
			fmt.Println("Sleeping")
			time.Sleep(4 * time.Second)
			importStats(repo, org, client, sem)
			return
		}
	}

	for _, s := range stats {
		for _, w := range s.Weeks {
			if *w.Commits == 0 {
				continue
			}

			var stat c.RepoStat
			var user c.User
			user.FromGhContrib(s.Author)

			user_id := findOrCreateUser(&user)
			// Does it exist?
			timeStr := fmt.Sprintf("%4d%02d%02d", w.Week.Time.Year(), w.Week.Time.Month(), w.Week.Time.Day())
			c.DB.Where("user_id = ? and repository_id = ? and to_char(week, 'YYYYMMDD') = ?",
				user_id,
				repo.Id,
				timeStr).First(&stat)

			if stat.Id == 0 || stat.Commits != int64(*w.Commits) {
				stat.UserId = user_id
				stat.RepositoryId = repo.Id
				stat.Week = w.Week.Time
				stat.Add = int64(*w.Additions)
				stat.Del = int64(*w.Deletions)
				stat.Commits = int64(*w.Commits)
				c.DB.Save(&stat)
			}
		}
	}

	<-*sem
}

func importPulls(repo *c.Repository, org *c.Organization, client *github.Client, page int) {
	// Import pulls for a given repo
	opt := &github.PullRequestListOptions{}
	opt.State = "all"
	opt.ListOptions = github.ListOptions{Page: page, PerPage: 100}

	pulls, response, err := client.PullRequests.List(org.Login, repo.Name, opt)
	fmt.Println("Getting pulls", repo.Name, page)
	fmt.Println(response)
	if err != nil {
		fmt.Println("Error with repo pulls: ", repo.Name)
		fmt.Println(err)
		return
	}

	// There are more pages, lets fetch the next one
	if response.NextPage > 0 {
		importPulls(repo, org, client, response.NextPage)
	}

	for _, gh_pull := range pulls {
		var pull c.Pull

		// Does the commit exist?
		c.DB.Where("number = ? and repository_id = ?", *gh_pull.Number, repo.Id).First(&pull)

		if pull.Id == 0 {
			pull.RepositoryId = repo.Id
			pull.Title = *gh_pull.Title
			pull.Body = getStr(gh_pull.Body)
			pull.Admin = *gh_pull.User.SiteAdmin
			pull.Number = int64(*gh_pull.Number)
			pull.GhCreatedAt.Time = *gh_pull.CreatedAt
			pull.GhCreatedAt.Valid = true

			var user c.User
			user.FromGhUser(gh_pull.User)

			pull.UserId = findOrCreateUser(&user)
		}

		// If the pull has not been updated go to the next one
		if pull.GhUpdatedAt.Time == *gh_pull.UpdatedAt {
			continue
		}

		pull.GhUpdatedAt.Time = *gh_pull.UpdatedAt
		pull.GhUpdatedAt.Valid = true

		if gh_pull.MergedAt != nil {
			pull.MergedAt.Time = *gh_pull.MergedAt
			pull.MergedAt.Valid = true
		}

		c.DB.Save(&pull)
	}
}
