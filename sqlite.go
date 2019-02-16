package main

import(
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

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
