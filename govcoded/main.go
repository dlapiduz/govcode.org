package main

import (
	"github.com/go-martini/martini"
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

	m.Group("/repos", func(r martini.Router) {
		r.Get("", ReposIndex)
		r.Get("/:id", ReposShow)
	})

	m.Group("/orgs", func(r martini.Router) {
		r.Get("", OrgsIndex)
		r.Get("/:id", OrgsShow)
	})

	m.Run()
}

func ReposIndex(r render.Render) {

	c.DB.SetLogger(log.New(os.Stdout, "\r\n", 0))

	c.DB.LogMode(true)

	var results []c.Repository
	rows := c.DB.Table("repositories")
	rows = rows.Select("repositories.*, organizations.name as organization_name")
	rows = rows.Joins("inner join organizations on organizations.id = repositories.organization_id")
	rows = rows.Limit(10)
	rows.Debug().Scan(&results)

	r.Header().Add("Access-Control-Allow-Origin", "*")

	r.JSON(200, results)
}

func ReposShow(r render.Render, params martini.Params) {
	var result c.Repository
	c.DB.Where("id = ?", params["id"]).First(&result)
	r.Header().Add("Access-Control-Allow-Origin", "*")

	r.JSON(200, result)
}

func OrgsIndex(r render.Render) {
	var results []c.Organization
	c.DB.Find(&results)
	r.Header().Add("Access-Control-Allow-Origin", "*")

	r.JSON(200, results)
}

func OrgsShow(r render.Render, params martini.Params) {
	var result c.Organization
	c.DB.Where("id = ?", params["id"]).First(&result)
	r.Header().Add("Access-Control-Allow-Origin", "*")

	r.JSON(200, result)
}
