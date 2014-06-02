package main

import (
	// "fmt"
	c "govcode.org/common"
)

func generateStats() {
	// Commit Count
	rows, err := c.DB.Table("commits").Select(`
    EXTRACT(month from date) as month,
    EXTRACT(year from date) as year,
    repository_id,
    organization_id,
    count(commits.id) as num_commits`).Joins(`
    inner join repositories on repositories.id =
    commits.repository_id
    `).Group("month, year, repository_id, organization_id").Rows()
	defer rows.Close()

	c.PanicOn(err)

	for rows.Next() {
		var month, year, repository_id, organization_id, num_commits int64
		rows.Scan(&month, &year, &repository_id, &organization_id, &num_commits)

		stat := c.CommitOrgStats{
			Month:          month,
			Year:           year,
			RepositoryId:   repository_id,
			OrganizationId: organization_id,
		}

		c.DB.Where(stat).FirstOrInit(&stat)

		stat.CommitCount = num_commits

		c.DB.Save(&stat)
	}

	// // New Pulls count
	// rows, err := c.DB.Table("pulls").Select(`
	//    EXTRACT(month from gh_created_at) as month,
	//    EXTRACT(year from gh_created_at) as year,
	//    repository_id,
	//    organization_id,
	//    count(pulls.id) as num_pulls`).Joins(`
	//    inner join repositories on repositories.id =
	//    pulls.repository_id
	//    `).Group("month, year, repository_id, organization_id").Rows()
	// defer rows.Close()

	// c.PanicOn(err)

	// for rows.Next() {
	// 	var month, year, repository_id, organization_id, num_commits int64
	// 	rows.Scan(&month,
	// 		&year,
	// 		&repository_id,
	// 		&organization_id,
	// 		&num_commits)

	// 	stat := c.CommitOrgStats{
	// 		Month:          month,
	// 		Year:           year,
	// 		RepositoryId:   repository_id,
	// 		OrganizationId: organization_id,
	// 	}

	// 	c.DB.Where(stat).FirstOrInit(&stat)

	// 	stat.CommitCount = num_commits

	// 	c.DB.Save(&stat)
	// }

}
