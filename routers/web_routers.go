/**
* Created by GoLand.
* User: nankai
* Date: 2019-07-25
* Time: 12:20
 */

package routers

import (
	"go-websocket/controllers/home"
	"go-websocket/controllers/systems"
	"go-websocket/controllers/user"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	router.LoadHTMLGlob("views/**/*")

	// user group
	userRouter := router.Group("/user")
	{
		userRouter.GET("/list", user.List)
		userRouter.GET("/online", user.Online)
		userRouter.POST("/sendMessage", user.SendMessage)
		userRouter.POST("/sendMessageAll", user.SendMessageAll)
	}

	// system
	systemRouter := router.Group("/system")
	{
		systemRouter.GET("/state", systems.Status)
	}

	// home
	homeRouter := router.Group("/home")
	{
		homeRouter.GET("/index", home.Index)
	}

	// router.POST("/user/online", user.Online)
}
