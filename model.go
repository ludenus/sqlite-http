package main

type Options struct {
	ListeningAddress string
	SqliteDbFile     string
}

type AgentDataSrcRecord struct {
	Id      int    `json:"id"`
	QaData  string `json:"qa_data"`
	Testrun int    `json:"testrun"`
	Stamp   int    `json:"stamp"`
}

type AgentNotificationRecord struct {
	Id           int `json:"id"`
	Notification string `json:"notification"`
}

type Info struct {
	GitBranch string `json:"gitBranch"`
	GitCommit string `json:"gitCommit"`
}