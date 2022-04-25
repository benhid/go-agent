package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

type agentExecRes struct {
	Message      string `json:"message"`
	Error        string `json:"error"`
	StdErr       string `json:"stderr"`
	StdOut       string `json:"stdout"`
	ExecDuration int64  `json:"exec_duration"`
	MemUsage     int64  `json:"mem_usage"`
}

func main() {
	if faasProcess := os.Getenv("fprocess"); faasProcess == "" {
		log.Fatal("Provide a valid process via fprocess environmental variable.")
		return
	}
	http.HandleFunc("/_/health", makeHealthHandler())
	http.HandleFunc("/", makeRunHandler())

	log.Info("Starting server at port :80")

	log.Fatal(http.ListenAndServe(":80", nil))
}

func makeHealthHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Agent is running"))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func makeRunHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.WithError(err).Error("Received invalid request")
			http.Error(w, "Decode failed", http.StatusInternalServerError)
			return
		}

		log.Info("Received run request")

		// Call job handler.
		execResult := handler(data)

		js, err := json.Marshal(execResult)
		if err != nil {
			log.WithError(err).Error("Decode failed")
			http.Error(w, "Decode failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("X-Duration-Seconds", fmt.Sprintf("%d", execResult.ExecDuration))
		// w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(js)
	}
}
