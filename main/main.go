package main

import (
	"flag"
	"log"
	"net/http"
	"os"


	"DBHS/config"
	"DBHS/docs" // Import generated docs
	"DBHS/middleware"

	_ "github.com/swaggo/swag"
	"github.com/MarceloPetrucio/go-scalar-api-reference"
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
	config.Router.PathPrefix("/reference").HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			// SpecURL: "https://generator3.swagger.io/openapi.json",// allow external URL or local path file
			SpecURL: "./docs/swagger.json", 
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Simple API",
			},
			DarkMode: true,
		})
		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(htmlContent))	
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
