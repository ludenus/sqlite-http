# sqlite-http


## rebuild & run
```bash
$ killall sqlite-http ; rm -f ./sqlite-http ; go build sqlite-http.go && ./sqlite-http -l :8008 -f ./sqlite.db &
```

## check
```bash
$ curl -i -X GET http://localhost:8008/qa

$ curl -i -X POST http://localhost:8008/qa -d "{ \"id\":-0, \"qa_data\":\"`whoami`@`hostname`\", \"testrun\":-1, \"stamp\":`date +%s` }"
```