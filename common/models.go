package common

import (
	"time"
)

type Organization struct {
	Id     int64
	Name   string `sql:"size(255)"`
	Login  string `sql:"size(255)"`
	Ignore bool

	Repositories []Repository

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Repository struct {
	Id          int64
	GhId        int64
	Name        string `sql:"size(255)"`
	Forks       int64
	Watchers    int64
	Stargazers  int64
	Size        int64
	OpenIssues  int64
	Description string
	Language    string

	OrganizationId int64
	Organization   Organization

	Commits []Commit
	Pulls   []Pull

	Ignore bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Commit struct {
	Id           int64
	Sha          string `sql:"size(255)"`
	Message      string
	Date         time.Time
	RepositoryId int64
	UserId       int64
}

type User struct {
	Id        int64
	GhId      int64
	Login     string `sql:"size(255)"`
	AvatarUrl string `sql:"size(255)"`

	Commits []Commit
}

type Pull struct {
	Id           int64
	RepositoryId int64
	Title        string
	Body         string
	Admin        bool
	Number       int64
	State        string
	UserId       int64

	MergedAt    time.Time
	GhCreatedAt time.Time
	GhUpdatedAt time.Time

	CreatedAt time.Time
}

type CommitOrgStats struct {
	Id              int64
	RepositoryId    int64
	OrganizationId  int64
	Month           int64
	Year            int64
	CommitCount     int64
	NewPullCount    int64
	ClosedPullCount int64
}
