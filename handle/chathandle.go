package handle

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Awaken1119/assistant-voice/db"
	"github.com/gin-gonic/gin"
)

func SaveChatHandler(c *gin.Context) {
	var chatData db.ChatData

	// 2. 解析前端发送的 JSON 数据到结构体
	if err := c.ShouldBindJSON(&chatData); err != nil {
		// 解析失败（如格式错误、字段不匹配）
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "解析对话数据失败：" + err.Error(),
		})
		return
	}

	//// 3. 解析成功，获取数据并处理（示例：打印数据）
	//// 对话组ID
	//println("对话ID：", chatData.ID)
	//// 对话标题
	//println("对话标题：", chatData.Title)
	//// 用户名
	//println("用户：", chatData.UserName)
	//// 消息列表（遍历）
	//for i, msg := range chatData.Messages {
	//	println("消息", i+1, "：角色=", msg.Role, "，内容=", msg.Content)
	//}

	// 4. 业务处理：保存到数据库（示例逻辑）
	err := db.SaveToDatabase(chatData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "保存失败：" + err.Error(),
		})
		return
	}

	// 5. 返回成功响应（包含后端生成的对话ID，可选）
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"chatId": chatData.ID, // 或返回数据库生成的ID
		},
		"message": "对话保存成功",
	})
}

func GetHistoryHandle(c *gin.Context) {

	userName := c.Query("username") // 对应 URL 中的 userId=123
	//pageStr := c.Query("page")
	//limitStr := c.Query("limit")
	chats, err := db.GetHistory(userName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "获取历史对话失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"chats": chats,
		},
	})
}

func GetChatHandle(c *gin.Context) {
	userName := c.Query("username")
	chatId := c.Query("chatId")
	chats, err := db.GetChat(userName, chatId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "获取历史对话失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":       chats.ID,
			"messages": chats.Messages,
		},
	})
}

func DeleteChatHandler(c *gin.Context) {
	username := c.Param("username")
	chatId := c.Param("chatId")

	pwd, _ := os.Getwd()
	filePath := fmt.Sprintf(pwd + "/historyfile/" + username + "/txtfile/" + chatId + ".txt")

	err := os.Remove(filePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "删除对话失败：" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除对话成功",
	})

}
