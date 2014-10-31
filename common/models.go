package common

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/lib/pq"
)

type Organization struct {
	Id     int64
	Name   string `sql:"size(255)"`
	Login  string `sql:"size(255)"`
	Ignore bool

	Repositories []Repository

	CreatedAt time.Time
	// UpdatedAt time.Time
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
	Description string `sql:"type:text;"`
	Language    string

	// Stats
	LastCommit  pq.NullTime
	LastPull    pq.NullTime
	CommitCount int64

	OrganizationId int64

	// Accessor Fields
	OrganizationLogin string `sql:-`
	DaysSincePull     int64  `sql:-`
	DaysSinceCommit   int64  `sql:-`

	// Related fields
	Commits      []Commit
	Pulls        []Pull
	Organization Organization
	RepoStat     []RepoStat

	Ignore bool

	CreatedAt time.Time
	UpdatedAt time.Time

	HelpWantedIssueCount int64
}

func (r Repository) TableName() string {
	return "repositories"
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
	Id          int64
	GhId        int64
	Login       string `sql:"size(255)"`
	AvatarUrl   string `sql:"size(255)"`
	CommitCount int64
	OrgList     string `sql:"size(255)"`

	Commits []Commit
}

type Pull struct {
	Id           int64
	RepositoryId int64
	Title        string `sql:"type:text;"`
	Body         string `sql:"type:text;"`
	Admin        bool
	Number       int64
	State        string
	UserId       int64

	MergedAt    pq.NullTime
	GhCreatedAt pq.NullTime
	GhUpdatedAt pq.NullTime

	CreatedAt time.Time
}

type Issue struct {
	Id           int64
	RepositoryId int64
	Number       int64
	Title        string `sql:"type:text;"`
	Body         string `sql:"type:text;"`
	Url          string
	Labels       string `sql:"type:text;"`
}

type RepoStat struct {
	Id           int64
	RepositoryId int64
	UserId       int64

	Week    time.Time
	Add     int64
	Del     int64
	Commits int64
}

type RepoAggStat struct {
	Id           int64
	RepositoryId int64

	Week    time.Time
	Add     int64
	Del     int64
	Commits int64
}

type CommitOrgStats struct {
	Id              int64
	Week            time.Time
	OrganizationId  int64
	CommitCount     int64
	NewPullCount    int64
	ClosedPullCount int64
}

func (u *User) FromGhUser(gh_user *github.User) {
	u.Login = *gh_user.Login
	u.GhId = int64(*gh_user.ID)
	u.AvatarUrl = *gh_user.AvatarURL
}

func (u *User) FromGhContrib(gh_contrib *github.Contributor) {
	u.Login = *gh_contrib.Login
	u.GhId = int64(*gh_contrib.ID)
	u.AvatarUrl = *gh_contrib.AvatarURL
}

func (i *Issue) HelpWanted() bool {
	labels := strings.Split(i.Labels, ",")
	for _, label := range labels {
		if matched, _ := regexp.MatchString("(?i)help.*?wanted|want.*?help|need.*?help", label); matched {
			return true
		}
	}
	return false
}
