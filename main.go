package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// ====================================== main
//
// https://www.sohamkamani.com/blog/2017/10/18/golang-adding-database-to-web-application/

func main() {

	opts := ParseArgs(os.Args[1:]) // must not include program name to parse successfully
	fmt.Println(opts)

	db = initDb(opts.SqliteDbFile)
	defer db.Close()

	http.HandleFunc("/qa", requestHandler)
	err := http.ListenAndServe(opts.ListeningAddress, nil)
	if err != nil {
		panic(err)
	}
}

// ====================================== http

func requestHandler(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case "GET":
		// w.Header().Set("Server", "QA")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	case "POST":
		decoder := json.NewDecoder(req.Body)
		var data AgentDataSrcRecord
		err := decoder.Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {

			// TODO: https://www.sohamkamani.com/blog/2017/10/18/golang-adding-database-to-web-application/

			response, err := json.Marshal(data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				w.Write(response)
			}
		}
	default:
		http.Error(w, "Only GET and POST requests are supported", http.StatusMethodNotAllowed)
	}

}
