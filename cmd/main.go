package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
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

	client := &http.Client{}
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

	// Display request and response details in the browser
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "Request:")
	fmt.Fprintln(w, string(requestDump))
	fmt.Fprintln(w, "\nResponse:")
	fmt.Fprintln(w, string(responseDump))
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
