package main

import (
	c "github.com/dlapiduz/govcode.org/common"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/render"
	"time"
)

func main() {
	m := App()
	m.Run()
}

func App() *martini.ClassicMartini {
	m := martini.Classic()

	m.Use(gzip.All())

	m.Use(render.Renderer(render.Options{
		Directory: "templates",
	}))

	m.Use(cors.Allow(&cors.Options{
		AllowAllOrigins: true,
	}))

	m.Group("/repos", func(r martini.Router) {
		r.Get("", ReposIndex)
		r.Get("/:name", ReposShow)
	})

	m.Group("/orgs", func(r martini.Router) {
		r.Get("", OrgsIndex)
		r.Get("/:id", OrgsShow)
	})

	m.Group("/users", func(r martini.Router) {
		r.Get("", UserIndex)
		r.Get("/:id", UserShow)
	})

	m.Get("/stats", StatsIndex)

	return m
}

func ReposIndex(r render.Render) {
	var results []c.Repository
	rows := c.DB.Table("repositories")
	rows = rows.Select(`repositories.*, organizations.login as organization_login,
		date_part('day', now() - last_pull) as days_since_pull,
		date_part('day', now() - last_commit) as days_since_commit
		`)
	rows = rows.Joins("inner join organizations on organizations.id = repositories.organization_id")
	rows.Scan(&results)

	r.JSON(200, results)
}

func ReposShow(r render.Render, params martini.Params) {

	var repo c.Repository
	c.DB.Where("name = ?", params["name"]).First(&repo)
	c.DB.Model(&repo).Related(&repo.Organization)
	c.DB.Model(&repo).Related(&repo.RepoStat)

	r.JSON(200, repo)
}

func OrgsIndex(r render.Render) {
	var results []c.Organization
	c.DB.Find(&results)

	r.JSON(200, results)
}

func OrgsShow(r render.Render, params martini.Params) {
	var result c.Organization
	c.DB.Where("id = ?", params["id"]).First(&result)

	r.JSON(200, result)
}

func UserIndex(r render.Render) {
	var results []c.User
	c.DB.Find(&results)

	r.JSON(200, results)
}

func UserShow(r render.Render, params martini.Params) {
	var result c.User
	c.DB.Where("id = ?", params["id"]).First(&result)

	r.JSON(200, result)
}

func StatsIndex(r render.Render) {
	// Get the repo counts per org
	type repoCount struct {
		OrganizationLogin string
		RepoCount         int64
	}
	var repo_counts []repoCount
	rows := c.DB.Table("repositories")
	rows = rows.Select(`organizations.login as organization_login, 
		count(repositories.name) as repo_count
		`)
	rows = rows.Joins("inner join organizations on organizations.id = repositories.organization_id")
	rows = rows.Group("organizations.login")
	rows = rows.Order("repo_count desc")
	rows.Scan(&repo_counts)

	// Get commit stats per org per month for the past year
	type orgStat struct {
		OrganizationLogin string
		Week              time.Time
		Month             string
		Add               int64
		Del               int64
		Commits           int64
	}
	var org_stats []orgStat
	rows = c.DB.Debug().Table("repo_stats")
	rows = rows.Select(`organizations.login as organization_login,
		min(repo_stats.week) as week,
		TO_CHAR(repo_stats.week, 'Mon YYYY') as month,
		sum(repo_stats.add) as add,
		sum(repo_stats.del) as del,
		sum(repo_stats.commits) as commits
	`)
	rows = rows.Joins(`inner join repositories on repositories.id = repo_stats.repository_id
		inner join organizations on organizations.id = repositories.organization_id
	`)
	rows = rows.Where("repo_stats.week > now()::date - 365")
	rows = rows.Group("organizations.login, month")
	rows = rows.Order("week")
	rows.Scan(&org_stats)

	var result map[string]interface{}
	result = make(map[string]interface{})

	result["repo_counts"] = repo_counts
	result["org_stats"] = org_stats

	r.JSON(200, result)
}
