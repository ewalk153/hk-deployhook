package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var apiKey string
var authKey string

func loadSite() {
	apiKey = os.Getenv("API_KEY")
}

func loadKey() {
	authKey = os.Getenv("AUTH_KEY")
}

func port() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return port
}

func NewRelicRequest(params *url.Values) {
	pForm := url.Values{}
	pForm.Set("deployment[app_name]", params.Get("app"))
	pForm.Set("deployment[description]", params.Get("git_log"))
	pForm.Set("deployment[revision]", params.Get("head"))
	pForm.Set("deployment[user]", params.Get("user"))

	req, err := http.NewRequest("POST", "https://api.newrelic.com/deployments.xml", strings.NewReader(pForm.Encode()))
	req.Header.Add("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println("Response", resp)
	if err != nil {
		log.Println("Err", err)
	}
}

func main() {
	loadSite()
	loadKey()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if r.RequestURI[1:] == authKey {
			NewRelicRequest(&r.Form)
		} else {
			log.Println("No post, auth does not match request:", r.RequestURI[1:])
		}
		log.Println(r.RequestURI, r.Form)
		fmt.Fprintln(w, "done")
	})
	fmt.Println("Listening on port", port())
	log.Fatal(http.ListenAndServe(":"+port(), nil))
}
