package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Response struct {
	RequestID int               `json:"request_id"`
	Responses map[string]string `json:"dc_responses"`
}

type Request struct {
	Id string `json:"id"`
	Version string `json:"version"`
	Timeout uint `json:"timeout_ms"`
	AdRequests []AdRequest `json:"ad_requests"`
}

type AdRequest struct {
	URL string `json:"url"`
	Body string `json:"body"`
	Headers map[string]string `json:"headers"`
	Method string `json:"method"`
	Timeout uint `json:"timeout_ms"`
	AdnetID uint `json:"adnet_id"`
}

var cannedResponse = Response{
			RequestID: 5,
			Responses: map[string]string{
				"mopub": "some response",
				"dfp":   "some other response",
			},
		}

func main() {
	http.HandleFunc("/mediate", func(w http.ResponseWriter, r *http.Request) {
		// validate request
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		requestBytes, err := ioutil.ReadAll(r.Body)
		defer func() {
			err := r.Body.Close()
			if err != nil {
				log.Printf("could not close request body: %v", err)
			}
		}()

		var validatedRequest AdRequest
		err = json.Unmarshal(requestBytes, &validatedRequest)
		if err != nil {
			log.Printf("couldn't marshal request: %v", err)
			return
		}

		if validatedRequest.URL == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, err = fmt.Fprint(w, "missing request URL")
			if err != nil {
				log.Printf("couldn't write response body: %v", err)
				return
			}
		}

		// return response
		bytes, err := json.Marshal(cannedResponse)
		if err != nil {
			log.Printf("error unmarshalling response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = fmt.Fprint(w, string(bytes))
		if err != nil {
			log.Printf("couldn't write response: %v", err)
		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("server exited with error: %v", err)
		os.Exit(1)
	}
}
