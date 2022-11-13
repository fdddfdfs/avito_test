package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx"
	_ "net/http"
)

var conn *pgx.Conn
var router *gin.Engine
var routerAddr string

func main() {
	err := connectToDatabase()
	if err != nil {
		panic(err.Error())
	}

	routerAddr = ":8080"

	router = gin.Default()
	router.Use(cors.Default())

	router.GET("/balance/:id", getBalance)
	router.PUT("/balance", addBalance)
	router.POST("/reservation", reserveMoney)
	router.PUT("/reservation/accept", acceptRevenue)
	router.PUT("/reservation/cancel", unreserveMoney)
	router.GET("/report/:year/:month", getReport)
	router.GET("/report/transaction/:userID/:page/:sort", getBalanceTransactions)

	//config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"http://localhost:8080", "http://localhost:80", "http://localhost:3200"}

	err = router.Run(routerAddr)
	if err != nil {
		panic(err.Error())
	}

}

func extractConfig() pgx.ConnConfig {
	var config pgx.ConnConfig

	config.Host = "pgdatabase"
	config.Port = 5432
	config.User = "root"
	config.Password = "root"
	config.Database = "avito_test"

	return config
}

func connectToDatabase() error {
	config := extractConfig()
	var err error
	conn, err = pgx.Connect(config)

	if err != nil {
		return err
	}

	return nil
}
