package main

import (
	"encoding/json"
	"net/http"
	// "os"
	"code.google.com/p/goauth2/oauth"
	"fmt"
	"github.com/google/go-github/github"
	c "govcode.org/common"
	"io/ioutil"
	"sync"
)

func getStr(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func findOrCreateUser(gh_user *github.User) int64 {
	var user c.User

	c.DB.Where("gh_id = ?", *gh_user.ID).First(&user)

	if user.Id == 0 {
		user.Login = *gh_user.Login
		user.AvatarUrl = getStr(gh_user.AvatarURL)
		user.GhId = int64(*gh_user.ID)
		c.DB.Save(&user)
	}

	return user.Id
}

func runImport() (err error) {
	// gh_key := os.Getenv("GH_KEY")

	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: "e20e82ed78f8d31be5f0ce5d875ccf62d72b8df8"},
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
			gh_org, _, err := client.Organizations.Get(org_name)
			// Do not panic, probably 404 error
			if err != nil {
				fmt.Println(err)
				wg.Done()
				return
			}

			var org c.Organization

			c.DB.Where("login = ?", gh_org.Login).First(&org)

			if org.Id == 0 {
				fmt.Println("New org")
				org.Login = *gh_org.Login
				if gh_org.Name != nil {
					org.Name = *gh_org.Name
				}
				org.Ignore = false

				c.DB.Save(&org)
				fmt.Println("Adding: ", org.Name)
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
	if err != nil {
		fmt.Println("Error with org: ", org.Login)
		return
	}

	// There are more pages, lets fetch the next one
	if response.NextPage > 0 {
		fmt.Println("Getting more", org.Login, page)
		go importRepos(org, client, response.NextPage)
	}

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

			// importCommits(&repo, org, client, 1)
			importPulls(&repo, org, client, 1)
		}

	}
}

func importCommits(repo *c.Repository, org *c.Organization, client *github.Client, page int) {
	// Import commits for a given repo
	opt := &github.CommitsListOptions{}
	opt.ListOptions = github.ListOptions{Page: page, PerPage: 100}

	var last_commit c.Commit
	c.DB.Where("repository_id = ?", repo.Id).Order("date desc").First(&last_commit)

	if last_commit.Id > 0 {
		opt.Since = last_commit.Date
	}

	commits, response, err := client.Repositories.ListCommits(org.Login, repo.Name, opt)
	if err != nil {
		fmt.Println("Error with repo: ", repo.Name)
		return
	}

	// There are more pages, lets fetch the next one
	if response.NextPage > 0 {
		fmt.Println("Getting more", org.Login, repo.Name, page)
		importCommits(repo, org, client, response.NextPage)
	}

	for _, gh_commit := range commits {
		var commit c.Commit

		// Does the commit exist?
		c.DB.Where("sha = ? and repository_id = ?", *gh_commit.SHA, repo.Id).First(&commit)

		if commit.Id == 0 {
			commit.Sha = *gh_commit.SHA
			commit.Message = getStr(gh_commit.Commit.Message)
			commit.Date = *gh_commit.Commit.Author.Date
			commit.RepositoryId = repo.Id
			commit.UserId = findOrCreateUser(gh_commit.Author)
			c.DB.Save(&commit)
		}
	}
}

func importPulls(repo *c.Repository, org *c.Organization, client *github.Client, page int) {
	// Import pulls for a given repo
	opt := &github.PullRequestListOptions{}
	opt.State = "all"
	opt.ListOptions = github.ListOptions{Page: page, PerPage: 100}

	pulls, response, err := client.PullRequests.List(org.Login, repo.Name, opt)
	if err != nil {
		fmt.Println("Error with repo pulls: ", repo.Name)
		fmt.Println(err)
		return
	}

	// There are more pages, lets fetch the next one
	if response.NextPage > 0 {
		fmt.Println("Getting more pulls", org.Login, repo.Name, page)
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
			pull.GhCreatedAt = *gh_pull.CreatedAt
			pull.UserId = findOrCreateUser(gh_pull.User)
		}
		pull.GhUpdatedAt = *gh_pull.UpdatedAt

		if gh_pull.MergedAt != nil {
			pull.MergedAt = *gh_pull.MergedAt
		}

		c.DB.Save(&pull)
	}
}
