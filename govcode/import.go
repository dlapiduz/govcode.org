package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"code.google.com/p/goauth2/oauth"
	c "github.com/dlapiduz/govcode.org/common"
	"github.com/google/go-github/github"
)

func runImport() (err error) {
	gh_key := os.Getenv("GH_KEY")

	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: gh_key},
	}

	client := github.NewClient(t.Client())

	orgs := getAccounts()
	importOrgs(orgs, client) // Also pulls the repos

	concurrency := 10

	sem := make(chan bool, concurrency)

	var repos []c.Repository
	rows := c.DB.Table("repositories")
	rows = rows.Select(`repositories.*, organizations.login as organization_login,
		coalesce(date_part('day', now() - last_pull), -1) as days_since_pull,
		coalesce(date_part('day', now() - last_commit), -1) as days_since_commit
		`)
	rows = rows.Joins("inner join organizations on organizations.id = repositories.organization_id")
	rows = rows.Order("updated_at")
	rows.Scan(&repos)

	for _, r := range repos {
		sem <- true
		if r.OrganizationLogin == "" {
			<-sem
			continue
		}

		if !r.Ignore {
			err := importStats(&r, r.OrganizationLogin, client)
			if err != nil {
				fmt.Println("There has been an error with stats")
				if strings.Contains(err.Error(), "404") {
					<-sem
					continue
				}
				c.PanicOn(err)
			}

			err = importPulls(&r, r.OrganizationLogin, client, 1)
			if err != nil {
				fmt.Println("There has been an error with stats")
				if strings.Contains(err.Error(), "404") {
					<-sem
					continue
				}
				c.PanicOn(err)
			}

		}
		c.DB.Save(&r)
		<-sem
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	return nil
}

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

	orgs = append(orgs, []string{"afrl", "arcticlcc", "arm-doe", "bbginnovate",
		"blue-button", "ca-cst-sii", "chaos", "cocomans", "commercegov", "cooperhewitt",
		"eeoc", "energyapps", "fccdata", "federal-aviation-administration", "globegit",
		"gopleader", "government-services", "govxteam", "greatsmokymountainsnationalpark",
		"hhsdigitalmediaapiplatform", "hhsidealab", "imdprojects", "innovation-toolkit",
		"ioos", "irsgov", "jbei", "kbase", "m-o-s-e-s", "measureauthoringtool", "missioncommand",
		"nasa-gibs", "nasa-rdt", "nasajpl", "nationalguard", "ncats", "ncbitools",
		"ncpp", "ncrn", "ndar", "neogeographytoolkit", "nersc", "ngageoint", "ngds",
		"nhanes", "niem", "nist-bws", "nist-ics-sec-tb", "nmml",
		"noaa-orr-erd", "nrel-cookbooks", "ozone-development", "petsc", "pm-master",
		"servir", "smithsonian", "sunpy", "usbr", "usdeptveteransaffairs", "usgcrp",
		"usgin-models", "usgs-astrogeology", "usgs-cida", "usgs-owi", "usgs-r", "usindianaffairs",
		"usnistgov", "usps", "vhainnovations", "virtual-world-framework", "visionworkbench",
		"wfmrda"}...)

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
			fmt.Println("Getting Org", org_name)

			if err != nil {
				fmt.Println("Error fetching orgs")
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
	fmt.Printf("Getting Repos for org Login %s at page %d\n", org.Login, page)
	if err != nil {
		fmt.Println("Error with org: ", org.Login)
		c.PanicOn(err)
	}

	// There are more pages, lets fetch the next one
	if response.NextPage > 0 {
		importRepos(org, client, response.NextPage)
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
				repo.Ignore = false
			}
			repo.Forks = int64(*r.ForksCount)
			repo.Watchers = int64(*r.WatchersCount)
			repo.Stargazers = int64(*r.StargazersCount)
			repo.Size = int64(*r.Size)
			repo.OpenIssues = int64(*r.OpenIssuesCount)
			repo.Language = getStr(r.Language)

			c.DB.Save(&repo)

		}
	}

}

func importStats(repo *c.Repository, org_login string, client *github.Client) error {
	stats, res, err := client.Repositories.ListContributorsStats(org_login, repo.Name)
	fmt.Println("Getting stats", repo.Name)
	fmt.Println(res)
	if err != nil {
		if res != nil && res.StatusCode == 202 {
			fmt.Println("Sleeping")
			time.Sleep(4 * time.Second)
			importStats(repo, org_login, client)
			return nil
		} else {
			return err
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

	return nil
}

func importPulls(repo *c.Repository, org_login string, client *github.Client, page int) error {
	// Import pulls for a given repo
	opt := &github.PullRequestListOptions{}
	opt.State = "all"
	opt.ListOptions = github.ListOptions{Page: page, PerPage: 100}

	pulls, response, err := client.PullRequests.List(org_login, repo.Name, opt)
	fmt.Printf("Getting pulls for repo %s at page %d\n", repo.Name, page)
	fmt.Println(response)
	if err != nil {
		return err
	}

	// There are more pages, lets fetch the next one
	if response.NextPage > 0 {
		importPulls(repo, org_login, client, response.NextPage)
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
			if gh_pull.CreatedAt != nil {
				pull.GhCreatedAt.Time = *gh_pull.CreatedAt
				pull.GhCreatedAt.Valid = true
			}

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
	return nil
}
