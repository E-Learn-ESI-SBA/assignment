package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"madaurus/dev/assignment/app/routes"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var Db *sql.DB

func ConnectDatabse() {
	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("USER"), os.Getenv("DB_NAME"), os.Getenv("PASSWORD"))
	db, errSql := sql.Open("postgres", psqlSetup)
	if errSql != nil {
		log.Fatal("err when conneting to db", errSql)
		panic(errSql)
	}
	Db = db
	fmt.Println("Success connecting to db")

}

func main() {

	configCors := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowFiles:      true,
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-CSRF-Token", "hx-request", "hx-current-url"},
		MaxAge:          12 * time.Hour,
	}
	// db connection
	ConnectDatabse()

	server := gin.New()

	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to madaurus assignments service"})
	})

	server.Use(cors.New(configCors))
	routes.AssignmentsRoute(server, Db)
	routes.SubmissionsRoute(server, Db)
	routes.FilesRoute(server, Db)

	// Set up GoMigrate
	driver, err := postgres.WithInstance(Db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://./db/migrations", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	m.Up()
	if err != nil {
		log.Fatal("JWT_SECRET not set")

	}

	// Listen and Server in 0.0.0.0:8080
	err = server.Run(":8080")
	if err != nil {
		log.Printf("Error starting server %v", err)
	}
}
