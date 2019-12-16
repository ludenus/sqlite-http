package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	_ "github.com/mattn/go-sqlite3"
)

// TODO: refactor error handling https://blog.golang.org/error-handling-and-go

func infoRequestHandler(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case "GET":
		response, err := json.Marshal(GitInfo)
		if err != nil {
			log.Println(httpError(w, err, http.StatusInternalServerError))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)

	default:
		log.Println(httpError(w, errors.New("Method is not supported"), http.StatusMethodNotAllowed))
	}

}

func notificationsRequestHandler(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case "GET":

		likes, err := parseParam("like", req)
		if err != nil {
			log.Println(httpError(w, err, http.StatusBadRequest))
			return
		}

		rows, err := selectNotificationsLike(likes[0])
		if err != nil {
			log.Println(httpError(w, err, http.StatusInternalServerError))
			return
		}

		defer rows.Close()

		notificationRecords := make([]AgentNotificationRecord, 0)

		for rows.Next() {
			var record AgentNotificationRecord

			err := rows.Scan(&record.Id, &record.Notification)
			if err != nil {
				log.Println(httpError(w, err, http.StatusInternalServerError))
				return
			}
			notificationRecords = append(notificationRecords, record)
		}

		err = rows.Err()
		if err != nil {
			log.Println(httpError(w, err, http.StatusInternalServerError))
			return
		}

		response, err := json.Marshal(notificationRecords)
		if err != nil {
			log.Println(httpError(w, err, http.StatusInternalServerError))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)

	default:
		log.Println(httpError(w, errors.New("Method is not supported"), http.StatusMethodNotAllowed))
	}

}

func dataRequestHandler(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case "GET":
		rows, err := selectTestDataAll()
		if err != nil {
			log.Println(httpError(w, err, http.StatusInternalServerError))
			return
		}

		defer rows.Close()

		dataRecords := make([]AgentDataSrcRecord, 0)

		for rows.Next() {
			var record AgentDataSrcRecord

			err := rows.Scan(&record.Id, &record.QaData, &record.Testrun, &record.Stamp, &record.BlobData)
			if err != nil {
				log.Println(httpError(w, err, http.StatusInternalServerError))
				return
			}
			dataRecords = append(dataRecords, record)
		}

		err = rows.Err()
		if err != nil {
			log.Println(httpError(w, err, http.StatusInternalServerError))
			return
		}

		response, err := json.Marshal(dataRecords)
		if err != nil {
			log.Println(httpError(w, err, http.StatusInternalServerError))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)

	case "POST":
		decoder := json.NewDecoder(req.Body)

		var dataToInsert AgentDataSrcRecord
		var selectedData AgentDataSrcRecord

		err := decoder.Decode(&dataToInsert)
		if err != nil {
			log.Println(httpError(w, err, http.StatusInternalServerError))
			return
		}

		res, err := insertTestData(dataToInsert)
		if err != nil {
			log.Println(httpError(w, err, http.StatusInternalServerError))
			return
		}

		// make sure data inserted
		id, err := res.LastInsertId()
		if err != nil {
			log.Println(httpError(w, err, http.StatusInternalServerError))
			return
		}

		rows, err := selectTestDataById(id)
		if err != nil {
			log.Println(httpError(w, err, http.StatusInternalServerError))
			return
		}

		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&selectedData.Id, &selectedData.QaData, &selectedData.Testrun, &selectedData.Stamp, &selectedData.BlobData)
			if err != nil {
				log.Println(httpError(w, err, http.StatusInternalServerError))
				return
			}

			response, err := json.Marshal(selectedData)
			if err != nil {
				log.Println(httpError(w, err, http.StatusInternalServerError))
				return
			}

			// only 1 record can be returned by id
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write(response)
		}

		err = rows.Err()
		if err != nil {
			log.Println(httpError(w, err, http.StatusInternalServerError))
			return
		}

	default:
		log.Println(httpError(w, errors.New("Method is not supported"), http.StatusMethodNotAllowed))
	}

}

func httpError(w http.ResponseWriter, err error, code int) error {
	http.Error(w, err.Error(), code)
	debug.PrintStack()
	return err
}

func parseParam(paramName string, req *http.Request) ([]string, error) {
	parsed, ok := req.URL.Query()[paramName]
	if !ok || len(parsed[0]) < 1 {
		err := fmt.Errorf("Query parameter '%s' is missing", paramName)
		return parsed, err
	}
	return parsed, nil
}
