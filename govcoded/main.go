package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/render"
	c "govcode.org/common"
	"log"
	"os"
)

func main() {
	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Directory: "templates",
	}))

	m.Use(cors.Allow(&cors.Options{
		AllowOrigins: []string{"*"},
	}))

	m.Group("/repos", func(r martini.Router) {
		r.Get("", ReposIndex)
		r.Get("/:id", ReposShow)
	})

	m.Group("/orgs", func(r martini.Router) {
		r.Get("", OrgsIndex)
		r.Get("/:id", OrgsShow)
	})

	m.Group("/users", func(r martini.Router) {
		r.Get("", UserIndex)
		r.Get("/:id", UserShow)
	})

	m.Run()
}

func ReposIndex(r render.Render) {

	c.DB.SetLogger(log.New(os.Stdout, "\r\n", 0))

	c.DB.LogMode(true)

	var results []c.Repository
	rows := c.DB.Table("repositories")
	rows = rows.Select("organizations.login as organization_login, repositories.*")
	rows = rows.Joins("inner join organizations on organizations.id = repositories.organization_id")
	rows.Scan(&results)

	r.JSON(200, results)
}

func ReposShow(r render.Render, params martini.Params) {
	var result c.Repository
	c.DB.Where("id = ?", params["id"]).First(&result)

	r.JSON(200, result)
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
