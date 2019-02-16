package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var (
	agentDataSrcTable      = "agent_data_src"
	agentNotificationTable = "agent_notifications"
)

func initDb(dbFile string) *sql.DB {

	db := openDbOrDie(createFileIfNotExists(dbFile))

	createTableDataSrc := fmt.Sprintf("CREATE TABLE  IF NOT EXISTS `%s` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `qa_data` VARCHAR(255) NULL, `testrun` INTEGER NULL,`stamp` INTEGER NULL)", agentDataSrcTable)
	createTableNotifications := fmt.Sprintf("CREATE TABLE  IF NOT EXISTS `%s` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `notification` VARCHAR(255) NULL)", agentNotificationTable)

	execSqlOrDie(db, createTableDataSrc)
	execSqlOrDie(db, createTableNotifications)

	return db
}

func insertTestData(data AgentDataSrcRecord) (sql.Result, error) {
	insertSql := fmt.Sprintf("INSERT INTO `%s` (`qa_data`, `testrun`, `stamp`) VALUES ($1,$2,$3)", agentDataSrcTable)
	return db.Exec(insertSql, data.QaData, data.Testrun, data.Stamp)
}

func selectTestData(id int64) (*sql.Rows, error) {
	selectSql := fmt.Sprintf("SELECT * FROM `%s` WHERE id=%d", agentDataSrcTable, id)
	return db.Query(selectSql)
}

func openDbOrDie(dbFile string) *sql.DB {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}
	return db
}

func execSqlOrDie(db *sql.DB, str string) sql.Result {
	res, err := db.Exec(str)
	if err != nil {
		panic(err)
	}
	return res
}