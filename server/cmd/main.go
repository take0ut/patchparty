package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/take0ut/patch-party/server/internal/config"
	"github.com/take0ut/patch-party/server/internal/models"
	"github.com/take0ut/patch-party/server/internal/routes"
)

func setupRoutes(app *config.Application) {
	routes.SetupDeviceRoutes(app)
	routes.SetupPatchRoutes(app)
	routes.SetupUserRoutes(app)
}

func main() {
	// Connect to the PostgreSQL database
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	app := &config.Application{
		Log:     log.New(os.Stdout, "patch-party ", log.LstdFlags),
		Router:  gin.Default(),
		Devices: &models.DeviceModel{DB: db},
		Users:   &models.UserModel{DB: db},
		Patches: &models.PatchModel{DB: db},
	}

	setupRoutes(app)

	// Run the server on port 8000
	if err := app.Router.Run(":8000"); err != nil {
		log.Fatal(err)
	}

	// Define a route handler for the GET / endpoint
	app.Router.GET("/", func(c *gin.Context) {
		// Send the webapp
		c.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})

}
