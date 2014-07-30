package main

import (
	"fmt"
	c "github.com/dlapiduz/govcode.org/common"
	"time"
)

func generateStats() {
	// Commit Count per Org
	fmt.Println("Generating commit count per org")

	rows, err := c.DB.Table("repo_stats").Select(`
	   week, organization_id, sum(commits) as commit_count`).Joins(`
	   inner join repositories
	   on repositories.id = repo_stats.repository_id
	   `).Group("week, repository_id, organization_id").Rows()

	c.PanicOn(err)

	for rows.Next() {
		var week time.Time
		var organization_id, commit_count int64
		rows.Scan(&week, &organization_id, &commit_count)

		stat := c.CommitOrgStats{
			Week:           week,
			OrganizationId: organization_id,
		}

		c.DB.Where(stat).FirstOrInit(&stat)

		stat.CommitCount = commit_count

		c.DB.Save(&stat)
	}

	rows.Close()

	// Update users commit count and org list
	fmt.Println("Generating commit count per user")

	rows, err = c.DB.Table("repo_stats").Select(`
	   user_id, array_agg(distinct o.login) as organization_list, sum(commits) as commit_count
	 `).Joins(`
	 	inner join repositories r
	 	on r.id = repo_stats.repository_id
	   inner join organizations o
	   on o.id = r.organization_id
	 `).Group("user_id").Rows()

	c.PanicOn(err)

	for rows.Next() {
		var user c.User
		var user_id, commit_count int64
		var org_list string
		rows.Scan(&user_id, &org_list, &commit_count)

		if user_id == 0 {
			continue
		}

		c.DB.Where("id = ?", user_id).First(&user)

		if commit_count != user.CommitCount {
			user.CommitCount = commit_count
			user.OrgList = org_list
			c.DB.Save(&user)
		}
	}

	rows.Close()

	// Update LastCommit
	fmt.Println("Updating LastCommit")

	rows, err = c.DB.Table("repo_stats").Select(`
    repository_id, max(week) as week
  `).Group("repository_id").Rows()

	c.PanicOn(err)

	for rows.Next() {
		var repository_id int64
		var week time.Time

		rows.Scan(&repository_id, &week)

		var repo c.Repository

		c.DB.Where("id = ?", repository_id).First(&repo)

		if repo.LastCommit.Time != week {
			c.DB.Model(&repo).Update("last_commit", week)
		}
	}

	rows.Close()

	// Update LastPull
	fmt.Println("Updating LastPull")

	rows, err = c.DB.Table("pulls").Select(`
	   repository_id, max(gh_created_at) as latest_pull
	 `).Group("repository_id").Rows()

	c.PanicOn(err)

	for rows.Next() {
		var repository_id int64
		var latest_pull time.Time

		rows.Scan(&repository_id, &latest_pull)

		var repo c.Repository

		c.DB.Where("id = ?", repository_id).First(&repo)

		if repo.LastPull.Time != latest_pull {
			repo.LastPull.Time = latest_pull
			repo.LastPull.Valid = true
			c.DB.Save(&repo)
		}
	}

	rows.Close()
}
