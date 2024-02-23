package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
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

	// Extract IP and headers from Via header
	viaHeader := resp.Header.Get("Via")
	viaEntries := strings.Split(viaHeader, ", ")
	ipAndHeaders := make([][]string, len(viaEntries))
	for i, entry := range viaEntries {
		ipAndHeaders[i] = strings.SplitN(entry, " ", 2)
	}

	// Display request and response details along with IP and headers in a table
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "<html><head><title>Proxy Hops</title></head><body>")
	fmt.Fprintln(w, "<h2>Proxy Hops:</h2>")
	fmt.Fprintln(w, "<table border=\"1\"><tr><th>IP</th><th>Headers Added</th></tr>")
	for _, entry := range ipAndHeaders {
		fmt.Fprintf(w, "<tr><td>%s</td><td>%s</td></tr>\n", entry[0], entry[1])
	}
	fmt.Fprintln(w, "</table>")
	fmt.Fprintln(w, "<h2>Request:</h2>")
	fmt.Fprintln(w, "<pre>")
	fmt.Fprintln(w, string(requestDump))
	fmt.Fprintln(w, "</pre>")
	fmt.Fprintln(w, "<h2>Response:</h2>")
	fmt.Fprintln(w, "<pre>")
	fmt.Fprintln(w, string(responseDump))
	fmt.Fprintln(w, "</pre>")
	fmt.Fprintln(w, "</body></html>")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
