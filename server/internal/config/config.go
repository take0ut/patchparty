package config

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/take0ut/patch-party/server/internal/models"
)

type Application struct {
	Log     *log.Logger
	Router  *gin.Engine
	Devices *models.DeviceModel
	Users   *models.UserModel
	Patches *models.PatchModel
}
