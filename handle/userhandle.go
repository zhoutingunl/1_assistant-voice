package handle

import (
	"net/http"

	"github.com/Awaken1119/assistant-voice/db/mysql"
	"github.com/Awaken1119/assistant-voice/util"

	"github.com/gin-gonic/gin"
)

const (
	pwd_sha = "*#@90"
)

// 定义与前端对应的响应结构体

func UserSignInHandler(c *gin.Context) {
	http.Redirect(c.Writer, c.Request, "/static/login.html", 302)
	return
}

func DoUserSignInHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	passwordHash := util.Sha1([]byte(pwd_sha + password))

	var user mysql.User
	db := mysql.InitDB()
	if err := db.Model(&mysql.User{}).Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err,
		})
	}

	if passwordHash == user.PasswordHash {
		c.JSON(http.StatusOK, gin.H{
			"success":      true,
			"message":      "登录成功",
			"token":        "aadasd",
			"refreshToken": "asdasd",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "密码错误",
		})
	}
	return
}

func UserRegisterHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	db := mysql.InitDB()
	var user mysql.User
	user.Username = username
	user.PasswordHash = util.Sha1([]byte(pwd_sha + password))

	if err := db.Model(mysql.User{}).Create(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "注册错误",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"message":      "注册成功",
		"token":        "aadasd",
		"refreshToken": "asdasd",
	})

}
