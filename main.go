package main

import (
	"net/http"

	"jwt-authen-golang-example/api"
	"jwt-authen-golang-example/service"
	"log"

	"io/ioutil"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const projectID = "smooth-delivery"

func main() {
	serviceAccount, err := ioutil.ReadFile("service-account.json")
	if err != nil {
		log.Fatal(err)
	}
	err = api.Init(api.Config{
		ServiceAccountJSON: serviceAccount,
		ProjectID:          projectID,
	})
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(
		middleware.Recover(),
		middleware.Secure(),
		middleware.Logger(),
		middleware.Gzip(),
		middleware.BodyLimit("2M"),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{
				"http://localhost:8080",
			},
			AllowHeaders: []string{
				echo.HeaderOrigin,
				echo.HeaderContentLength,
				echo.HeaderAcceptEncoding,
				echo.HeaderContentType,
				echo.HeaderAuthorization,
			},
			AllowMethods: []string{
				echo.GET,
				echo.POST,
			},
			MaxAge: 3600,
		}),
	)

	// Health check
	e.GET("/_ah/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Register services
	service.Auth(e.Group("/auth"))

	e.Logger.Fatal(e.Start(":9000"))
	// e.Run(standard.WithConfig(engine.Config{
	// 	Address:      ":9000",
	// 	ReadTimeout:  30 * time.Second,
	// 	WriteTimeout: 30 * time.Second,
	// }))
}
