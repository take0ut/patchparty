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

func SetupUserRoutes(app *config.Application) {
	app.Router.POST("/login", func(c *gin.Context) {
		var user models.User
		c.BindJSON(&user)
		user, err := app.Users.GetUserByEmail(user.Email)
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
		if auth.VerifyPassword([]byte(user.Password), user.Password) {
			app.Users.SetLastLogin(user.ID)
			c.Header("patchparty_session", auth.CreateSession(app.Users, user.ID))
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

	app.Router.POST("/logout", func(c *gin.Context) {
		var user models.User
		c.BindJSON(&user)
		err := app.Users.SetUserSession(user.ID, "")
		if err != nil {
			c.JSON(200, gin.H{
				"message": "User failed to log out!",
				"error":   err,
			})
		}
		c.Cookie("patchparty_session=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;")
		c.JSON(200, gin.H{
			"message": "User logged out!",
		})
	})

	app.Router.POST("/users", func(c *gin.Context) {
		var user models.User
		c.BindJSON(&user)
		user.Password = auth.HashAndSalt([]byte(user.Password))
		// Insert user into database
		user, err := app.Users.CreateUser(user.UserName, user.Email, user.Password)
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
		session := auth.CreateSession(app.Users, user.ID)

		c.Cookie(session)
		c.JSON(200, gin.H{
			"message": "User created!",
			"user":    user,
		})
	})
}
