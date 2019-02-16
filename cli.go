package main

import (
	"os"
	"flag"
)

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
