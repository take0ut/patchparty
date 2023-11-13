package routes

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/take0ut/patch-party/server/internal/config"
	"github.com/take0ut/patch-party/server/internal/models"
)

func SetupPatchRoutes(app *config.Application) {
	app.Router.GET("/patches/:id", func(c *gin.Context) {
		id := c.Param("id")
		patch, err := app.Patches.GetPatchById(id)
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
		c.JSON(200, gin.H{
			"message": "Patch found!",
			"patch":   patch,
		})
		return
	})

	app.Router.POST("/patches", func(c *gin.Context) {
		var patch models.Patch
		c.BindJSON(&patch)
		patch, err := app.Patches.CreatePatch(patch)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(200, gin.H{
					"message": "Patch not created.",
				})
				return
			} else {
				log.Fatal(err)
			}
		}
		c.JSON(200, gin.H{
			"message": "Patch created!",
			"patch":   patch,
		})
	})
}
