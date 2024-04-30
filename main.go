package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/knadh/koanf"

	"madaurus/dev/assignment/app/shared"

	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
_ "github.com/golang-migrate/migrate/v4/source/file"
)

var Db *sql.DB

func ConnectDatabse(k *koanf.Koanf) {
	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", k.String("HOST"), k.Int("PORT"), k.String("USER"), k.String("DB_NAME"), k.String("PASSWORD"))
	fmt.Println(psqlSetup)
	db, errSql := sql.Open("postgres", psqlSetup)
	if errSql != nil {
		log.Fatal("err when conneting to db", errSql)
		panic(errSql)
	}
	Db = db
	fmt.Println("Success connecting to db")

}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r
}

func main() {
	k := shared.GetSecrets()
	r := setupRouter()

	// db connection 
	ConnectDatabse(k)

	// Set up GoMigrate
	driver, err := postgres.WithInstance(Db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
    m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
    m.Up() 

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
