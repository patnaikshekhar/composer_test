package main

import (
	"os"
	"sample/internal/db"
	"sample/internal/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	dbType := os.Getenv("DB_TYPE")
	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")

	conn, err := db.New(dbType, dbConnectionString)
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", conn)
			return next(c)
		}
	})

	e.GET("/api/v1/healthz", routes.Healthz)

	e.Static("/", "ui/dist")

	err = e.Start(":9081")
	if err != nil {
		e.Logger.Fatal(err)
	}
}
