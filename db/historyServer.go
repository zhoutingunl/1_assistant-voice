package db

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Awaken1119/assistant-voice/db/mysql"
)

type Message struct {
	ID       int64  `json:"id"`      // 消息唯一ID（时间戳）
	Content  string `json:"content"` // 消息内容
	Role     string `json:"role"`    // 角色：user/assistant
	Time     string `json:"time"`    // 消息时间（如"09:45"）
	Loading  bool   `json:"loading"` // 是否为加载状态
	AudioUrl string `json:"audioUrl"`
}

// ChatData 对话组结构体（对应前端的 chatData）
type ChatData struct {
	ID        int64     `json:"id"`        // 对话组ID
	Title     string    `json:"title"`     // 对话标题
	Timestamp string    `json:"timestamp"` // 对话日期（如"2025-10-20"）
	Messages  []Message `json:"messages"`  // 消息列表（嵌套结构体数组）
	UserName  string    `json:"username"`  // 用户名（注意前端是userId，后端结构体字段首字母大写）
}

func SaveToDatabase(chatData ChatData) error {
	db := mysql.InitDB()
	var user mysql.User
	var chat mysql.ChatHistory

	pwd, _ := os.Getwd()
	filePath := fmt.Sprintf(pwd + "/historyfile/" + chatData.UserName + "/txtfile/" + strconv.FormatInt(chatData.ID, 10) + ".txt")

	// 1️⃣ 先确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	fileData, _ := json.Marshal(chatData)
	err := os.WriteFile(filePath, fileData, 0644)
	if err != nil {
		return err
	}

	db.Model(mysql.User{}).Where("username = ?", chatData.UserName).Find(&user)

	chat.UserName = chatData.UserName
	chat.UserID = user.ID
	chat.Filepath = pwd + "/historyfile/" + chatData.UserName

	if err = db.Model(mysql.ChatHistory{}).FirstOrCreate(&chat).Error; err != nil {
		return err
	}
	return nil
}

func GetHistory(username string) ([]ChatData, error) {
	db := mysql.InitDB()
	var user mysql.User
	var chats []mysql.ChatHistory
	var history ChatData
	var historyData []ChatData

	if err := db.Model(mysql.User{}).Where("username = ?", username).Find(&user).Error; err != nil {
		return nil, err
	}
	if err := db.Model(mysql.ChatHistory{}).Where("user_id = ?", user.ID).Find(&chats).Error; err != nil {
		return nil, err
	}
	for _, chatHistory := range chats {
		paths, err := searchTxt(chatHistory.Filepath)
		if err != nil {
			return nil, err
		}
		for _, path := range paths {
			chat, err := os.ReadFile(path)
			if err != nil {
				return nil, err
			}
			_ = json.Unmarshal(chat, &history)
			historyData = append(historyData, history)
		}

	}
	return historyData, nil
}
func searchTxt(folder string) ([]string, error) {
	var txtFiles []string
	// 遍历文件夹
	err := filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 判断是否是 .txt 文件
		if !d.IsDir() && filepath.Ext(d.Name()) == ".txt" {
			txtFiles = append(txtFiles, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("遍历文件夹出错:", err)
		return nil, err
	}

	//// 输出所有 txt 文件路径
	//fmt.Println("找到的 .txt 文件如下：")
	//for _, file := range txtFiles {
	//	fmt.Println(file)
	//}
	return txtFiles, nil
}

func GetChat(username string, chatId string) (ChatData, error) {
	db := mysql.InitDB()
	var user mysql.User
	var chats mysql.ChatHistory
	var history ChatData

	if err := db.Model(mysql.User{}).Where("username = ?", username).Find(&user).Error; err != nil {
		return ChatData{}, err
	}
	if err := db.Model(mysql.ChatHistory{}).Where("user_id = ?", user.ID).Find(&chats).Error; err != nil {
		return ChatData{}, err
	}

	data, err := os.ReadFile(fmt.Sprintf(chats.Filepath + "/txtfile/" + chatId + ".txt"))
	if err != nil {
		return ChatData{}, err
	}
	_ = json.Unmarshal(data, &history)

	return history, nil
}
