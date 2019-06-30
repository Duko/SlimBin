package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var requests = make([]string, 0)

func apiResponse(w http.ResponseWriter, r *http.Request) {
	requestDump, _ := httputil.DumpRequest(r, true)
	decodedValue, _ := url.QueryUnescape(string(requestDump))
	fmt.Println(decodedValue)
	requests = append(requests, decodedValue)
}

func inspectResponse(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index").Parse(`
		<!DOCTYPE html>
		<html lang="en">
			<head>
				<title>Requests</title>
				<style>
					.request {
						border: 1px solid grey;
						padding: 10px;
					}

					.request:not(:last-child) {
						border-bottom: none;
					}
				</style>
			</head>
			<body>
				{{range .}}
					<div class="request">
						<pre>{{.}}</pre>
					</div>
				{{end}}
			</body>
		</html>
	`))

	tmpl.Execute(w, requests)
}

func getEnv(key, fallback string) string {
	if env, ok := os.LookupEnv(key); ok {
		return env
	}
	return fallback
}

func main() {
	finish := make(chan bool)

	api := http.NewServeMux()
	api.HandleFunc("/", apiResponse)
	go func() {
		if err := http.ListenAndServe(":"+getEnv("LISTEN_PORT", "80"), api); err != nil {
			panic(err)
		}
	}()

	inspect := http.NewServeMux()
	inspect.HandleFunc("/", inspectResponse)
	go func() {
		if err := http.ListenAndServe(":"+getEnv("INSPECT_PORT", "8080"), inspect); err != nil {
			panic(err)
		}
	}()

	<-finish
}
