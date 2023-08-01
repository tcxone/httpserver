package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s %s", "INFO", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
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
				nameInfo, nameErr := os.Stat(path + "/" + name)
				if nameErr != nil {
				}
				if nameInfo.IsDir() {
					fmt.Fprintf(w, "<li><a href=\"/%s\">%s/</a></li>\n", filepath.Join(path, name), name)
				} else {
					fmt.Fprintf(w, "<li><a href=\"/%s\">%s</a></li>\n", filepath.Join(path, name), name)
				}
			}
			fmt.Fprintln(w, "</ul><hr></body></html>")
		} else {
			mimeType := mime.TypeByExtension(filepath.Ext(path))

			if mimeType == "" {
				mimeType = "application/octet-stream"
			}

			w.Header().Set("Content-Type", mimeType)

			content, _ := ioutil.ReadFile(path)
			fmt.Fprintf(w, "%s", content)
		}
	})
	var port int
	fmt.Print("input server port >>>")
	fmt.Scanln(&port)
	/*if len(os.Args)>1{
		port = os.Args[len(os.Args)-1]
	}*/
	log.Println("[INFO] Start server on 0.0.0.0:" + fmt.Sprintf("%d", port))
	http.ListenAndServe(fmt.Sprintf(":%d", port), loggingMiddleware(http.DefaultServeMux))
}
