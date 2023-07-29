package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[%s] %s %s %s", "INFO", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s", "INFO", "Done in", time.Since(start))
	})
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := "./" + r.URL.Path
		info, err := os.Stat(path)
		if err != nil {
			http.Error(w, "404 Not Found", 404)
			return
		}
		if info.IsDir() {
			files, _ := ioutil.ReadDir(path)
			var list []string
			for _, file := range files {
				list = append(list, file.Name())
			}
			fmt.Fprintln(w, "<html lang='en'><head><meta charset='utf-8'><title>Directory listing for "+r.URL.Path+"</title></head><body><h1>Directory listing for "+r.URL.Path+"</h1><hr><ul>")
			for _, name := range list {
				fmt.Fprintf(w, "<li><a href=\"/%s\">%s</a></li>\n", filepath.Join(path, name), name)
			}
			fmt.Fprintln(w, "</ul><hr></body></html>")
		} else {
			content, _ := ioutil.ReadFile(path)
			fmt.Fprintf(w, "%s", content)
		}
	})
	log.Println("Start server on 0.0.0.0:8088")
	http.ListenAndServe(":8088", loggingMiddleware(http.DefaultServeMux))
}
