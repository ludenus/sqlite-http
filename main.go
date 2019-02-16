package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
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
	log.Println(opts)

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

		var dataToInsert AgentDataSrcRecord
		var selectedData AgentDataSrcRecord

		err := decoder.Decode(&dataToInsert)
		if err != nil {
			reportError500(w, err)
		} else {
			res, err := insertTestData(dataToInsert)
			if err != nil {
				reportError500(w, err)
			} else {
				id, err := res.LastInsertId()
				if err != nil {
					reportError500(w, err)
				} else {
					rows, err := selectTestData(id)
					if err != nil {
						reportError500(w, err)
					} else {
						defer rows.Close()

						for rows.Next() {
							err := rows.Scan(&selectedData.Id, &selectedData.QaData, &selectedData.Testrun, &selectedData.Stamp)
							if err != nil {
								reportError500(w, err)
							} else {
								response, err := json.Marshal(selectedData)
								if err != nil {
									reportError500(w, err)
								} else {
									w.Header().Set("Content-Type", "application/json")
									w.WriteHeader(http.StatusCreated)
									w.Write(response)
								}
							}
						}

						err := rows.Err()
						if err != nil {
							reportError500(w, err)
						}
					}
				}
			}
		}
	default:
		reportError(w, errors.New("Only GET and POST requests are supported"), http.StatusMethodNotAllowed)
	}

}

func reportError500(w http.ResponseWriter, err error) {
	reportError(w, err, http.StatusInternalServerError)
}

func reportError(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
	log.Println(err)
}
