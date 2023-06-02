package routers

import (
	"github.com/bryant-rh/srew/cmd/server/routers/plugin"
	"github.com/bryant-rh/srew/cmd/server/routers/user"
	"github.com/bryant-rh/srew/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func NewRooter(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	user.UserRouter(v1)
	v1.Use(middleware.JWTAuthMiddleware())
	{
		plugin.PluginRouter(v1)

	}

}
