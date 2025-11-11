package handle

import (
	"fmt"
	"net/http"

	"github.com/Awaken1119/assistant-voice/util"
	"github.com/gin-gonic/gin"
)

type VoiceRequest struct {
	Command   string                   `json:"command" binding:"required"`
	Context   []map[string]interface{} `json:"context"`
	Timestamp string                   `json:"timestamp" binding:"required"`
}

func ProcessVoiceHandler(c *gin.Context) {
	// 1. 声明一个结构体变量用于接收数据
	var req VoiceRequest

	// 2. 解析前端发送的JSON数据到结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		// 解析失败（如字段缺失、格式错误）
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误：" + err.Error(),
		})
		return
	}

	var reqString string
	for i, con := range req.Context {
		if i == len(req.Context)-2 {
			reqString = fmt.Sprintf(reqString + "\n" + "<question>" + con["content"].(string) + "</question>")
		} else {
			reqString = fmt.Sprintf(reqString + "\n" + "历史记录：" + con["role"].(string) + ": " + con["content"].(string))
		}
	}
	fmt.Printf("req:%s", reqString)
	output, err := util.LlmRun(reqString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "返回数据错误：" + err.Error(),
		})
		return
	}

	tts, err := util.VoiceTts(output)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true, // 标识请求处理成功
		"message": "成功",
		"data": gin.H{ // data字段：存放具体业务数据
			"message":  output, // 助手要显示的回复内容
			"audioUrl": tts,
		},
	})
	return
}
