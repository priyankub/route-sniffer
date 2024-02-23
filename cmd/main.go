package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "URL parameter is required", http.StatusBadRequest)
		return
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Error creating request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	resp, err := client.Do(request)
	if err != nil {
		http.Error(w, "Error sending request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Dump request and response details
	requestDump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		http.Error(w, "Error dumping request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	responseDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		http.Error(w, "Error dumping response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Split headers by new lines
	requestHeaders := strings.Split(string(requestDump), "\n")
	responseHeaders := strings.Split(string(responseDump), "\n")

	// Display request and response details in the browser
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1>Request</h1>")
	fmt.Fprintf(w, "<table border=\"1\"><tr><th>IP</th><th>Headers Added</th></tr>")
	for _, header := range requestHeaders {
		fmt.Fprintf(w, "<tr><td>%s</td><td></td></tr>", header)
	}
	fmt.Fprintf(w, "</table>")

	fmt.Fprintf(w, "<h1>Response</h1>")
	fmt.Fprintf(w, "<table border=\"1\"><tr><th>IP</th><th>Headers Added</th></tr>")
	for _, header := range responseHeaders {
		fmt.Fprintf(w, "<tr><td></td><td>%s</td></tr>", header)
	}
	fmt.Fprintf(w, "</table>")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
