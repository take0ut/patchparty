package routes

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/take0ut/patch-party/server/internal/auth"
	"github.com/take0ut/patch-party/server/internal/config"
	"github.com/take0ut/patch-party/server/internal/models"
)

func setupUserRoutes(app *config.Application) {
	app.Router.POST("/login", func(c *gin.Context) {
		var user models.User
		c.BindJSON(&user)

		row := db.QueryRow(queries.GetUserByEmail, user.Email)
		err := row.Scan(&user.ID, &user.UserName, &user.Email, &user.Password)
		if err != nil {
			log.Fatal(err)
		}
		if auth.VerifyPassword([]byte(user.Password), user.Password) {
			app.Users.SetLastLogin(user.ID)
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
