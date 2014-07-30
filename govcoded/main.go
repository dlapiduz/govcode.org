package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/render"
	c "govcode.org/common"
)

func main() {
	m := App()
	m.Run()
}

func App() *martini.ClassicMartini {
	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Directory: "templates",
	}))

	m.Use(cors.Allow(&cors.Options{
		AllowOrigins: []string{"*"},
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
