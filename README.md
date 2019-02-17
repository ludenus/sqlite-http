# sqlite-http
provides web interface for inserting qa data into sqlite db

## build
```bash
$ go build -v github.com/ludenus/sqlite-http && ./sqlite-http -l :8008 -f ./sqlite.db &
```

## build small binary
```bash
$ GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v github.com/ludenus/sqlite-http && upx --ultra-brute ./sqlite-http
```

## run
```bash
$ ./sqlite-http -l :8008 -f ./sqlite.db &
```

## check status
```bash
$ curl -i -X GET http://localhost:8008/qa
```

## insert qa data into sqlite db
```bash
$ curl -i -X POST http://localhost:8008/qa -d "{ \"id\":-0, \"qa_data\":\"`whoami`@`hostname`\", \"testrun\":-1, \"stamp\":`date +%s` }"
```