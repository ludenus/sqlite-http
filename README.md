# sqlite-http
provides web interface for inserting qa data into sqlite db

## build
```bash
$ ./build.sh
```

## build small binary
```bash
$ ./build.sh upx
```

## run
```bash
$ ./sqlite-http -l :8008 -f ./sqlite.db &
```

## check info
```bash
$ curl -X GET http://localhost:8008/info

{"gitBranch":"master","gitCommit":"bbc25078"}
```

## insert qa data into sqlite db
```bash
$ curl -v -X POST http://localhost:8008/data -d "{ \"id\":-0, \"qa_data\":\"`whoami`@`hostname`\", \"testrun\":-1, \"stamp\":`date +%s`, \"blob_data\":\"`whoami`@`hostname`\" }"

{"id":5,"qa_data":"andrey@qa-pc","testrun":-1,"stamp":1576503117,"blob_data":"andrey@qa-pc"}
```

## select qa data from sqlite db
```bash
$ curl -X GET http://localhost:8008/data

[{"id":1,"qa_data":"andrey@qa-pc","testrun":-1,"stamp":1576502460,"blob_data":""}]
```

## select notifications from sqlite db
```bash
$ curl -X GET http://localhost:8008/notification?like=sqlNotification%25

[{"Id":1,"Notification":"sqlNotification1550691478665"}]
```