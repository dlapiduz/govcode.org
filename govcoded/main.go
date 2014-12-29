package main

import (
	"net/http"
	"strings"
	"time"

	c "github.com/dlapiduz/govcode.org/common"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/render"
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

	m.Get("", Index)

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
	m.Get("/issues", IssuesIndex)

	return m
}

func Index(r render.Render) {
	type homeStats struct {
		RepoCount  int
		IssueCount int
		UserCount  int
		IssueLangs []string
	}

	stats := homeStats{}

	c.DB.Model(c.Repository{}).Count(&stats.RepoCount)
	c.DB.Model(c.Issue{}).Count(&stats.IssueCount)
	c.DB.Model(c.User{}).Count(&stats.UserCount)

	rows := c.DB.Raw(`select distinct language
		from issues
		inner join repositories on repositories.id = issues.repository_id
		where language != ''
		order by language
	`)

	rows.Pluck("language", &stats.IssueLangs)

	r.JSON(200, stats)
}

func ReposIndex(r render.Render, req *http.Request) {
	qs := req.URL.Query()
	perPage := ForceStoInt(qs.Get("perPage"))

	var results []c.Repository
	rows := c.DB.Table("repositories")
	rows = rows.Select(`repositories.*, organizations.login as organization_login,
		coalesce(date_part('day', now() - last_pull), -1) as days_since_pull,
		coalesce(date_part('day', now() - last_commit), -1) as days_since_commit
		`)
	rows = rows.Joins("inner join organizations on organizations.id = repositories.organization_id")
	if perPage > 0 {
		rows = rows.Limit(perPage)
	}
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
	rows = c.DB.Table("repo_stats")
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

func IssuesIndex(r render.Render, req *http.Request) {
	qs := req.URL.Query()

	opts := struct {
		perPage       int
		page          int
		repoId        int
		orgId         int
		languages     string
		organizations string
		state         string
		label         string
	}{
		ForceStoInt(qs.Get("perPage")),
		ForceStoInt(qs.Get("page")),
		ForceStoInt(qs.Get("repoId")),
		ForceStoInt(qs.Get("orgId")),
		qs.Get("languages"),
		qs.Get("organizations"),
		qs.Get("state"),
		qs.Get("label"),
	}

	var issues []c.Issue
	rows := c.DB.Table("issues").Select(`issues.*, 
		organizations.login as organization_login,
		repositories.language as language,
		repositories.name as repository_name
	`)
	rows = rows.Joins(`inner join repositories on repositories.id = issues.repository_id
		inner join organizations on organizations.id = repositories.organization_id
	`)

	if opts.repoId > 0 {
		rows = rows.Where("repositories.id = ?", opts.repoId)
	}
	if opts.orgId > 0 {
		rows = rows.Where("organizations.id = ?", opts.orgId)
	}
	if opts.languages != "" {
		langs := strings.Split(opts.languages, ",")
		rows = rows.Where("repositories.language = ANY (ARRAY [?])", langs)
	}
	if opts.organizations != "" {
		orgs := strings.Split(opts.organizations, ",")
		rows = rows.Where("organizations.login = ANY (ARRAY [?])", orgs)
	}
	if opts.label != "" {
		label := "%" + opts.label + "%"
		rows = rows.Where("issues.labels LIKE ?", label)
	}

	if opts.state == "all" || opts.state == "closed" {
		rows = rows.Where("issues.state = ?", opts.state)
	} else {
		rows = rows.Where("issues.state = ?", "open")
	}

	if opts.perPage == 0 || opts.perPage > 100 {
		opts.perPage = 100
	}

	rows = rows.Limit(opts.perPage)

	if opts.page > 0 {
		rows = rows.Offset(opts.page * opts.perPage)
	}

	rows = rows.Order("issues.gh_updated_at desc")

	rows.Scan(&issues)

	r.JSON(200, issues)
}
