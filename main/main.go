package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"DBHS/config"
	"DBHS/docs" // Import generated docs
	"DBHS/middleware"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/swaggo/swag"
)

// @title DBHS API
// @version 1.0
// @description API for DBHS application
// @termsOfService http://swagger.io/terms/

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.

// @BasePath /api
func main() {
	// ---- http server ---- //

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Programmatically set swagger info
	docs.SwaggerInfo.Title = "DBHS API"
	docs.SwaggerInfo.Description = "API for DBHS application"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8000"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	config.Init(infoLog, errorLog)
	defer config.CloseDB()

	addr := flag.String("addr", "0.0.0.0:8000", "HTTP network address")
	flag.Parse()

	defineURLs()

	// Initialize Swagger documentation
	// Add swagger endpoints to the router
	config.Router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The URL pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	// Serve ReDoc UI
	config.Router.HandleFunc("/redoc", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(".", "docs", "redoc.html"))
	})

	// Serve Scalar UI
	config.Router.HandleFunc("/scalar", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(".", "docs", "scalar.html"))
	})

	// Directly serve swagger.json for Scalar UI
	config.Router.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeFile(w, r, filepath.Join(".", "docs", "swagger.json"))
	})

	handler := middleware.LimitMiddleware(config.Router)

	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  handler,
	}

	infoLog.Printf("starting server on :%s", *addr)
	infoLog.Printf("Swagger UI available at http://localhost%s/swagger/index.html", *addr)
	infoLog.Printf("ReDoc UI available at http://localhost%s/redoc", *addr)
	infoLog.Printf("Scalar UI available at http://localhost%s/scalar", *addr)

	err := server.ListenAndServe()
	errorLog.Fatal(err)
}
