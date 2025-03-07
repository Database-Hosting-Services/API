package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"DBHS/config"
	_ "DBHS/docs"
)

//@title My API
//@version 1.0
//@description This is my API using Swagger with Mux.
//@host localhost:8080
//@BasePath /api/v1

func main() {
	// ---- http server ---- //

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	config.Init(infoLog, errorLog)
	defer config.CloseDB()

	addr := flag.String("addr", ":8000", "HTTP network address")
	flag.Parse()

	defineURLs()

	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  config.Mux,
	}

	infoLog.Printf("starting server on :%s", *addr)
	err := server.ListenAndServe()
	errorLog.Fatal(err)
}
