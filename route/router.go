package route

import (
	"github.com/Awaken1119/assistant-voice/handle"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	//router.LoadHTMLGlob("static/*")
	router.Static("/static", "./static")
	router.GET("/user/login", handle.UserSignInHandler)
	router.POST("/user/login", handle.DoUserSignInHandler)
	router.POST("user/register", handle.UserRegisterHandler)

	router.POST("/voice/process", handle.ProcessVoiceHandler)

	router.POST("/chat/save", handle.SaveChatHandler)
	router.GET("/chat/history", handle.GetHistoryHandle)
	router.GET("/chat/get", handle.GetChatHandle)
	router.DELETE("/chat/delete/:username/:chatId", handle.DeleteChatHandler)

	return router
}
