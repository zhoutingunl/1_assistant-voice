package appAction

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	myredis "github.com/Awaken1119/assistant-voice/cache/redis"
	"github.com/gomodule/redigo/redis"
	"github.com/mitchellh/go-ps"
)

type OpenAppAction struct {
	Action  string `json:"action"`
	AppName string `json:"appname"`
}

type SearchWebAction struct {
	Action   string `json:"action"`
	Question string `json:"question"`
}

type WriteTextAction struct {
	Action   string `json:"action"`
	Content  string `json:"content"`
	FilePath string `json:"filepath"`
}

type OpenFileAction struct {
	Action   string `json:"action"`
	AppName  string `json:"appname"`
	FilePath string `json:"filepath"`
}

type PlayMusicAction struct {
	Action    string `json:"action"`
	MusicName string `json:"musicname"`
}

type CloseAppAction struct {
	Action  string `json:"action"`
	AppName string `json:"appname"`
}

var appPIDList = make(map[string]string)

func OpenApp(action string) (string, error) {
	var cmd *exec.Cmd
	var openApps OpenAppAction

	rConn := myredis.InitRedis().Get()
	defer rConn.Close()

	_ = json.Unmarshal([]byte(action), &openApps)
	appPath, _ := redis.String(rConn.Do("GET", openApps.AppName))

	cmd = exec.Command(appPath)

	err := cmd.Start()
	if err != nil {
		fmt.Println("❌ 打开失败:", err)
		return "", err
	}
	p, err := ps.FindProcess(cmd.Process.Pid)
	fmt.Println(p.Executable())
	appPIDList[openApps.AppName] = p.Executable()
	return fmt.Sprintf(openApps.AppName + "已打开"), nil
}
func WriteText(action string) (string, error) {
	var writeText WriteTextAction
	_ = json.Unmarshal([]byte(action), &writeText)

	pwd, _ := os.Getwd()
	filePath := fmt.Sprintf(pwd + "/historyfile" + "/file" + writeText.FilePath)

	// 1️⃣ 先确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	fileData, _ := json.Marshal(writeText.Content)
	err := os.WriteFile(filePath, fileData, 0644)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("已写入文件：" + filePath), nil
}

// SearchWeb 打开默认浏览器并搜索关键字
func SearchWeb(query string) (string, error) {
	var cmd *exec.Cmd
	var Search SearchWebAction
	_ = json.Unmarshal([]byte(query), &Search)
	searchURL := "https://www.google.com/search?q=" + url.QueryEscape(Search.Question)

	cmd = exec.Command("cmd", "/C", "start", "", searchURL)

	//return fmt.Sprintf("浏览器已打开，已搜索:" + Search.Question + "。操作已完成，结束对话"), cmd.Start()
	return fmt.Sprintf("操作已完成，结束对话"), cmd.Start()
}

func OpenFile(action string) (string, error) {
	var cmd *exec.Cmd
	var openFile OpenFileAction

	rConn := myredis.InitRedis().Get()
	defer rConn.Close()

	_ = json.Unmarshal([]byte(action), &openFile)
	appPath, _ := redis.String(rConn.Do("GET", openFile.AppName))

	// 启动 Word 并打开指定文件
	cmd = exec.Command(appPath, openFile.FilePath)

	err := cmd.Start()
	if err != nil {
		fmt.Println("❌ 打开失败:", err)
		return "", err
	}

	return fmt.Sprintf("已用" + openFile.AppName + "打开:" + openFile.FilePath), nil
}

func PlayMusic(action string) (string, error) {
	var cmd *exec.Cmd
	var playMusic PlayMusicAction

	rConn := myredis.InitRedis().Get()
	defer rConn.Close()

	_ = json.Unmarshal([]byte(action), &playMusic)
	appPath, _ := redis.String(rConn.Do("GET", "酷狗音乐"))
	cmd = exec.Command(appPath)
	err := cmd.Start()
	if err != nil {
		fmt.Println("❌ 打开失败:", err)
		return "", err
	}

	MouseClick(playMusic.MusicName)
	return fmt.Sprintf("已播放音乐" + playMusic.MusicName), nil
}

func CloseApp(action string) (string, error) {
	var cmd *exec.Cmd
	var cliseApps CloseAppAction

	rConn := myredis.InitRedis().Get()
	defer rConn.Close()

	_ = json.Unmarshal([]byte(action), &cliseApps)

	appPID := appPIDList[cliseApps.AppName]
	fmt.Println(appPID)
	cmd = exec.Command("taskkill", "/IM", appPID, "/F")
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}

	return fmt.Sprintf(cliseApps.AppName + "已关闭," + "结束对话"), nil
}
