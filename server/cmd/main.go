package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"

	"github.com/take0ut/patch-party/server/internal/auth"
	"github.com/take0ut/patch-party/server/internal/config"
	queries "github.com/take0ut/patch-party/server/internal/db"
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
	aoor.GET("/", func(c *gin.Context) {
		// Send the webapp
		c.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})

	r.GET("/devices/:id", func(c *gin.Context) {
		device, err := models.DeviceModel.GetDeviceById(id)
		if err != nil || device == nil {
			if err == sql.ErrNoRows {
				c.JSON(200, gin.H{
					"message": "Device not found!",
				})
				return
			} else {
				log.Fatal(err)
			}
		}
		c.JSON(200, gin.H{
			"message": "Device found!",
			"device":  device,
		})
	})

	r.GET("/patches/:id", func(c *gin.Context) {
		var patch Patch
		id := c.Param("id")
		row := db.QueryRow(queries.GetPatchById, id)
		err := row.Scan(&patch.ID, &patch.Name, &patch.Description, &patch.DeviceID, &patch.UserID, &patch.Downloads, &patch.URI)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(200, gin.H{
					"message": "User not found!",
				})
				return
			} else {
				log.Fatal(err)
			}
		}
	})

	r.POST("/patches", func(c *gin.Context) {
		var patch Patch
		c.BindJSON(&patch)
		// Insert patch into database
		db.Query(queries.CreatePatch, patch.Name, patch.Description, patch.DeviceID, patch.UserID, patch.Downloads, patch.URI)
		// Return patch object
		c.JSON(200, gin.H{
			"message": "Patch created!",
			"patch":   patch,
		})
	})

	r.POST("/login", func(c *gin.Context) {
		var user User
		c.BindJSON(&user)
		row := db.QueryRow(queries.GetUserByEmail, user.Email)
		err := row.Scan(&user.ID, &user.UserName, &user.Email, &user.Password)
		if err != nil {
			log.Fatal(err)
		}
		if auth.VerifyPassword([]byte(user.Password), user.Password) {
			db.Query(queries.SetLastLogin, user.ID)
			c.JSON(200, gin.H{
				"message": "User logged in!",
				"user":    user,
			})
		} else {
			c.JSON(200, gin.H{
				"message": "User not logged in!",
				"user":    user,
			})
		}
	})

	r.POST("/logout", func(c *gin.Context) {
		var user User
		c.BindJSON(&user)
		db.Query(queries.SetSession, user.ID, nil)
		c.JSON(200, gin.H{
			"message": "User logged out!",
			"user":    user,
		})
	})

	r.POST("/users", func(c *gin.Context) {
		var user User
		c.BindJSON(&user)
		user.Password = auth.HashAndSalt([]byte(user.Password))
		// Insert user into database
		err := db.QueryRow(queries.CreateUser, user.UserName, user.Email, user.Password).Scan(&user)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(200, gin.H{
					"message": "User not created!",
				})
				return
			} else if err.(*pq.Error).Code == "23505" {
				c.JSON(401, gin.H{
					"message": "User already exists with this email address or username.",
				})
			} else {
				log.Fatal(err)
			}
		}
		// Return user object
		c.Header("patchparty_session", auth.CreateSession(db, user.ID))
		c.JSON(200, gin.H{
			"message": "User created!",
			"user":    user,
		})
	})

}
