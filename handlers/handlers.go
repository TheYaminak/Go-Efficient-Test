package handlers

import "net/http"

func StartServers() {
	http.HandleFunc("/process_sequential", SequentialHandler)
	go func() {
		http.ListenAndServe(":8080", nil)
	}()

	http.HandleFunc("/process_concurrent", ConcurrentHandler)
	http.ListenAndServe(":8081", nil)
}
