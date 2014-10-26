package common

import (
	"fmt"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// Connection string parameters for Postgres - http://godoc.org/github.com/lib/pq, if you are using another
// database refer to the relevant driver's documentation.

// * dbname - The name of the database to connect to
// * user - The user to sign in as
// * password - The user's password
// * host - The host to connect to. Values that start with / are for unix domain sockets.
//   (default is localhost)
// * port - The port to bind to. (default is 5432)
// * sslmode - Whether or not to use SSL (default is require, this is not the default for libpq)
//   Valid SSL modes:
//    * disable - No SSL
//    * require - Always SSL (skip verification)
//    * verify-full - Always SSL (require verification)

var DB gorm.DB

func init() {
	var err error

	cmd := os.Args[0]
	if strings.Contains(cmd, "test") {
		// We are doing testing!
		DB, err = gorm.Open("sqlite3", ":memory:")

		fmt.Println("TEST")
		DB.AutoMigrate(Organization{})
		DB.AutoMigrate(Repository{})
		DB.AutoMigrate(Commit{})
		DB.AutoMigrate(User{})
		DB.AutoMigrate(Pull{})
		DB.AutoMigrate(CommitOrgStats{})
		DB.AutoMigrate(RepoStat{})
	} else {
		DB, err = gorm.Open("postgres", os.Getenv("PG_CONN_STR"))
		// DB.LogMode(true)
		DB.DB().SetMaxOpenConns(10)
	}
	if err != nil {
		panic(fmt.Sprintf("Got error when connect database, the error is '%v'", err))
	}
}
