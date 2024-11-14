// @title Banking API
// @version 1.0
// @description This is a simple banking API.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"Accounts/api"
	database "Accounts/internal/datab"
	"Accounts/notifications"
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "Accounts/docs"

	_ "github.com/lib/pq"
)

var queries *database.Queries

func InitDB() *sql.DB {
	
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Read environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Construct the connection string
	connStr := "user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " host=" + dbHost + " port=" + dbPort + " sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Database is unreachable:", err)
	}

	queries = database.New(db) // Initialize queries with the database connection
	return db
}

func initializeKafka(notificationService *notifications.NotificationService, db *database.Queries) {
	brokerList := []string{os.Getenv("KAFKA_BROKER_LIST")}
	api.Initializekafka(brokerList)
	ctx:= context.TODO()
	go notifications.StartConsumer(brokerList, "transaction-notifications", *notificationService, db, ctx)
}

func main() {
	// Initialize the database and set up the router
	db := InitDB()
	err := godotenv.Load()
	if err != nil {
    	log.Fatalf("Error loading .env file")
	}

	emailHost := os.Getenv("EMAIL_HOST")
	emailFrom := os.Getenv("EMAIL_FROM")
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	port := os.Getenv("EMAIL_PORT")
	twilioSID:= os.Getenv("TWILIO_ACCOUNT_SID")
	twilioAuthToken:= os.Getenv("TWILIO_AUTH_TOKEN")
	twilioPhoneNumber:= os.Getenv("TWILIO_FROM_NUMBER")
	emailPort, err := strconv.Atoi(port)
	if err != nil {
		log.Printf("Failed to convert emailPort to string!")
		return
	}
	notificationService := notifications.NewNotificationService(emailFrom, emailPassword, emailHost, emailPort, twilioSID, twilioAuthToken, twilioPhoneNumber)
	initializeKafka(notificationService, queries)

	router := api.SetupRoutes(queries, db)

	// Initialize Swagger endpoint on the parent router
	router.Handle("/swagger/*", httpSwagger.WrapHandler)
	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}