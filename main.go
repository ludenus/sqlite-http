package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// ====================================== main

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	opts := ParseArgs(os.Args[1:]) // must not include program name to parse successfully

	db = initDb(opts.SqliteDbFile)
	defer db.Close()

	http.HandleFunc("/qa", requestHandler)

	log.Println("listeining on: " + opts.ListeningAddress)

	err := http.ListenAndServe(opts.ListeningAddress, nil)
	if err != nil {
		panic(err)
	}
}

// TODO: refactor error handling https://blog.golang.org/error-handling-and-go

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
			log.Println(err)
			reportError(w, err, http.StatusInternalServerError)
			return
		}

		res, err := insertTestData(dataToInsert)
		if err != nil {
			log.Println(err)
			reportError(w, err, http.StatusInternalServerError)
			return
		}

		id, err := res.LastInsertId()
		if err != nil {
			log.Println(err)
			reportError(w, err, http.StatusInternalServerError)
			return
		}

		rows, err := selectTestData(id)
		if err != nil {
			log.Println(err)
			reportError(w, err, http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&selectedData.Id, &selectedData.QaData, &selectedData.Testrun, &selectedData.Stamp)
			if err != nil {
				log.Println(err)
				reportError(w, err, http.StatusInternalServerError)
				return
			}

			response, err := json.Marshal(selectedData)
			if err != nil {
				log.Println(err)
				reportError(w, err, http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write(response)
		}

		err = rows.Err()
		if err != nil {
			log.Println(err)
			reportError(w, err, http.StatusInternalServerError)
			return
		}

	default:
		reportError(w, errors.New("Only GET and POST requests are supported"), http.StatusMethodNotAllowed)
	}

}

func reportError(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
	debug.PrintStack()
}
