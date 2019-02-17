package main

import (
	"flag"
	"os"
)

func ParseArgs(arguments []string) Options {
	var options = Options{
		ListeningAddress: fromEnvVar("SQLITE_HTTP_LISTENING_ADDRESS", ":8008"),
		SqliteDbFile:     fromEnvVar("SQLITE_HTTP_DB_FILE", "sqlite.db"),
	}

	fs := flag.NewFlagSet("sqlite-http", flag.ExitOnError)

	fs.StringVar(&options.ListeningAddress, "l", options.ListeningAddress, "listening address ip:port")
	fs.StringVar(&options.SqliteDbFile, "f", options.SqliteDbFile, "sqlite db file name")

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
