package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

var GitBranch string
var GitCommit string

var GitInfo = Info{
	GitBranch: GitBranch,
	GitCommit: GitCommit,
}

// ====================================== main

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("GitInfo: %s", GitInfo)

	opts := ParseArgs(os.Args[1:]) // must not include program name to parse successfully

	db = initDb(opts.SqliteDbFile)
	defer db.Close()

	http.HandleFunc("/info", infoRequestHandler)
	http.HandleFunc("/data", dataRequestHandler)
	http.HandleFunc("/notifications", notificationsRequestHandler)

	log.Println("listeining on: " + opts.ListeningAddress)

	err := http.ListenAndServe(opts.ListeningAddress, nil)
	if err != nil {
		panic(err)
	}
}
