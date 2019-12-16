package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var (
	agentDataSrcTable      = "agent_data_src"
	agentNotificationTable = "agent_notifications"
)

func initDb(dbFile string) *sql.DB {

	f := createFileIfNotExists(dbFile)
	log.Println("opening sqlite db: " + f)

	db := openDbOrDie(f)

	createTableDataSrc := fmt.Sprintf("CREATE TABLE  IF NOT EXISTS `%s` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `qa_data` VARCHAR(255) NULL, `testrun` INTEGER NULL,`stamp` INTEGER NULL, `blob_data` VARCHAR(255) NULL)", agentDataSrcTable)
	createTableNotifications := fmt.Sprintf("CREATE TABLE  IF NOT EXISTS `%s` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `notification` VARCHAR(255) NULL)", agentNotificationTable)

	execSqlOrDie(db, createTableDataSrc)
	execSqlOrDie(db, createTableNotifications)

	return db
}

func insertTestData(data AgentDataSrcRecord) (sql.Result, error) {
	sql := fmt.Sprintf("INSERT INTO `%s` (`qa_data`, `testrun`, `stamp`, `blob_data`) VALUES ($1,$2,$3,$4)", agentDataSrcTable)
	return db.Exec(sql, data.QaData, data.Testrun, data.Stamp, data.BlobData)
}

func selectTestDataById(id int64) (*sql.Rows, error) {
	sql := fmt.Sprintf("SELECT * FROM `%s` WHERE id=%d", agentDataSrcTable, id)
	return db.Query(sql)
}

func selectTestDataAll() (*sql.Rows, error) {
	sql := fmt.Sprintf("SELECT * FROM `%s`", agentDataSrcTable)
	return db.Query(sql)
}

func selectNotificationById(id int64) (*sql.Rows, error) {
	sql := fmt.Sprintf("SELECT * FROM `%s` WHERE id=%d", agentNotificationTable, id)
	return db.Query(sql)
}

func selectNotificationsLike(str string) (*sql.Rows, error) {
	sql := fmt.Sprintf("SELECT * FROM `%s` WHERE notification LIKE $1", agentNotificationTable)
	return db.Query(sql, str)
}

func selectNotificationsAll() (*sql.Rows, error) {
	sql := fmt.Sprintf("SELECT * FROM `%s`", agentNotificationTable)
	return db.Query(sql)
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
