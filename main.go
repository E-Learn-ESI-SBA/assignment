package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/knadh/koanf"

	"madaurus/dev/assignment/app/routes"
	"madaurus/dev/assignment/app/shared"

	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var Db *sql.DB

func ConnectDatabse(k *koanf.Koanf) {
	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", k.String("HOST"), k.Int("PORT"), k.String("USER"), k.String("DB_NAME"), k.String("PASSWORD"))
	db, errSql := sql.Open("postgres", psqlSetup)
	if errSql != nil {
		log.Fatal("err when conneting to db", errSql)
		panic(errSql)
	}
	Db = db
	fmt.Println("Success connecting to db")

}

func main() {

	k := shared.GetSecrets()

	// db connection
	ConnectDatabse(k)

	server := gin.New()

	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to madaurus assignments service"})
	})

	routes.AssignmentsRoute(server, Db)
	routes.SubmissionsRoute(server, Db)

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
	err = os.Setenv("JWT_SECRET", k.String("JWT_SECRET"))
	if err != nil {
		log.Fatal("JWT_SECRET not set")

	}

	// Listen and Server in 0.0.0.0:8080
	err = server.Run(":8080")
	if err != nil {
		log.Printf("Error starting server %v", err)
	}
}
