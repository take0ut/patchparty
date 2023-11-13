package routes

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/take0ut/patch-party/server/internal/config"
	"github.com/take0ut/patch-party/server/internal/models"
)

func SetupDeviceRoutes(app *config.Application) {
	app.Router.GET("/devices/:id", func(c *gin.Context) {
		var device models.Device
		id := c.Param("id")
		device, err := app.Devices.GetDeviceById(id)
		if err != nil {
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
}
