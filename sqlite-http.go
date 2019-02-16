package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Options struct {
	ListeningAddress string
	SqliteDbFile     string
}

type AgentDataSrcRecord struct {
	Id int `json:"id"`
	QaData string `json:"qa_data"`
	Testrun int `json:"testrun"`
	Stamp int `json:"stamp"`
}

type AgentNotificationRecord struct {
	Id int
	Notification string
}

// ====================================== main

func main() {

	opts := ParseArgs(os.Args[1:]) // must not include program name to parse successfully
	fmt.Println(opts)

	db := initDb(opts.SqliteDbFile)
	defer db.Close()
	
	http.HandleFunc("/qa", requestHandler)
  	http.ListenAndServe(opts.ListeningAddress, nil)

}

// ====================================== http

func requestHandler(w http.ResponseWriter, req *http.Request) {

	// if r.URL.Path != "/" {
    //     http.Error(w, "404 not found.", http.StatusNotFound)
    //     return
    // }
 
    switch req.Method {
    case "GET":     
		// w.Header().Set("Server", "QA")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
    case "POST":
        decoder := json.NewDecoder(req.Body)
		var data AgentDataSrcRecord
		err := decoder.Decode(&data)
		if err != nil {
			panic(err)
		}
		fmt.Println(data)
		w.WriteHeader(201)
    default:
        fmt.Fprintf(w, "ERROR: only GET and POST methods are supported.")
	}
	
	
  }


// ====================================== db


func initDb(dbFile string) *sql.DB {

	db := openDbOrDie(createFileIfNotExists(dbFile))

	createTableDataSrc := "CREATE TABLE  IF NOT EXISTS `agent_data_src` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `qa_data` VARCHAR(255) NULL, `testrun` INTEGER NULL,`stamp` INTEGER NULL)"
	createTableNotifications := "CREATE TABLE  IF NOT EXISTS `agent_notifications` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `notification` VARCHAR(255) NULL)"

	execSqlOrDie(db, createTableDataSrc)
	execSqlOrDie(db, createTableNotifications)

	return db
}

func openDbOrDie(dbFile string) *sql.DB {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}
	return db
}

func execSqlOrDie(db *sql.DB, str string) {
	_, err := db.Exec(str)
	if err != nil {
		panic(err)
	}
}

// ====================================== file

func createFileIfNotExists(filename string) string {

	if fileExists(filename) {
		fmt.Println("file exists: " + filename)
	} else {
		fmt.Println("creating file: " + filename)

		os.MkdirAll(filepath.Dir(filename), 0755)
		os.Create(filename)
	}

	return filename
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// ======================================= args

func ParseArgs(arguments []string) Options {
	var options = Options{
		ListeningAddress: fromEnvVar("SQLITE_HTTP_LISTENING_ADDRESS", ":8008"),
		SqliteDbFile:     fromEnvVar("SQLITE_HTTP_DB_FILE", "sqlite.db"),
	}

	fs := flag.NewFlagSet("main", flag.ExitOnError)

	fs.StringVar(&options.ListeningAddress, "listening-address", options.ListeningAddress, "listen on ip:port")
	fs.StringVar(&options.ListeningAddress, "l", options.ListeningAddress, "listen on ip:port")

	fs.StringVar(&options.SqliteDbFile, "sqlite-db-file", options.SqliteDbFile, "clients description")
	fs.StringVar(&options.SqliteDbFile, "f", options.SqliteDbFile, "clients description")

	fs.Parse(arguments)
	return options
}

func fromEnvVar(envVarName string, defaultValue string) string {
	res := defaultValue
	fromEnv := os.Getenv(envVarName)
	if fromEnv != "" {
		res = fromEnv
	}
	return res
}
