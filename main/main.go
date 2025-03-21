package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"DBHS/config"
	"DBHS/middleware"
)

func main() {
	// ---- http server ---- //

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	config.Init(infoLog, errorLog)
	defer config.CloseDB()

	addr := flag.String("addr", "0.0.0.0:8000", "HTTP network address")
	flag.Parse()

	defineURLs()

	handler := middleware.LimitMiddleware(config.Router)

	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  handler,
	}

	infoLog.Printf("starting server on :%s", *addr)

	err := server.ListenAndServe()
	errorLog.Fatal(err)
}
