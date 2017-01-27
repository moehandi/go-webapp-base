package logrequest

import (
	"fmt"
	"net/http"
	"time"
	"github.com/fatih/color"
)

// Handler will log the HTTP requests
func Handler(next http.Handler) http.Handler {
	//return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Println(time.Now().Format("2006-01-02 03:04:05 PM"), r.RemoteAddr, r.Method, r.URL)
	//	next.ServeHTTP(w, r)
	//})

	color.Set(color.FgCyan, color.Bold)
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		//res := NewMyResponseWriter(w)
		next.ServeHTTP(w, r)
		//res.status = 200
		end := time.Now()

		fmt.Println("[LOG]", time.Now().Format("2006/01/02 - 03:04:05 PM"), r.RemoteAddr, r.Method, r.URL, end.Sub(start))
		//fmt.Println("[LOG]", time.Now().Format("2006/01/02 - 03:04:05 PM"), r.RemoteAddr, r.Method, r.URL, end.Sub(start), res.status)
	}
	return http.HandlerFunc(fn)
}

type MyResponseWriter struct {
	http.ResponseWriter
	status int
}

func NewMyResponseWriter(w http.ResponseWriter) *MyResponseWriter {
	// Default the status code to 200
	return &MyResponseWriter{w, http.StatusOK}
}

// Give a way to get the status
func (w MyResponseWriter) Status() int {
	return w.status
}

// Satisfy the http.ResponseWriter interface
func (w MyResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w MyResponseWriter) Write(data []byte) (int, error) {
	return w.ResponseWriter.Write(data)
}

func (w MyResponseWriter) WriteHeader(statusCode int) {
	// Store the status code
	w.status = statusCode

	// Write the status code onward.
	w.ResponseWriter.WriteHeader(statusCode)
}
