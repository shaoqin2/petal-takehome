package main

import (
	"encoding/json"
	"flag"
	"fmt"
	SHOUTCLOUD "github.com/RICHO/GOSHOUT"
	"gopkg.in/retry.v1"
	"log"
	"net/http"
)

var port = flag.Int("port", 8090, "the port to listen on")
var exponentialBackoff = retry.LimitCount(3, retry.Exponential{MaxDelay: 10})

type reverseParams struct {
	Data string `json:"data"`
}

func reverse(w http.ResponseWriter, req *http.Request) {
	if contentType := req.Header.Get("Content-Type"); contentType != "" {
		if contentType != "application/json" {
			http.Error(w, "Content-Type header is not application/json", http.StatusUnsupportedMediaType)
			return
		}
	}

	var data reverseParams
	// Note this is dangerous in a realistic production environment
	// we will be decoding as many bytes as client sent us. To protect against malicious users, we probably should set MaxBytesReader
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, "bad data input format, see api documentation for correct request body format", http.StatusBadRequest)
		return
	}

	for a := retry.Start(exponentialBackoff, nil); a.Next(); {
		// a timeout built is into this client library, though VERY janky, it has hard coded 20s and it is recreating HTTP client unnecessarily...
		// for the sake of a take home, we won't create another timeout system around this, but any http client should have a timeout to avoid overwhelming the server with hanging requests
		upcased, err := SHOUTCLOUD.UPCASE(ReverseString(data.Data))
		if err != nil {
			log.Printf("[%d] failed to contact shoutcloud.io: %s\n", http.StatusInternalServerError, err.Error())
			if a.More() {
				continue
			} else {
				http.Error(w, "server failed to query required upcase service, please reach out to admins", http.StatusInternalServerError)
				return
			}
		}
		log.Printf("[%d] upcased user input %s to %s", http.StatusOK, data.Data, upcased)
		fmt.Fprintf(w, upcased)
		return
	}
	// we should not reach here because either the last retry is successful and we return 200 or it didn't we return 500. But have this here just for safety
	http.Error(w, "server failed", http.StatusInternalServerError)
}

func main() {
	flag.Parse()
	// for such a simple application, just use global default handler
	// realistically in a production application, you'd want a struct with its own package / file for handlers

	// for takehome we use a simplistic go raw http server over a framework like echo. In production we might want a full REST framework
	http.HandleFunc("/v1", reverse)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatalf("server failed: %s\n", err.Error())
	}
}
