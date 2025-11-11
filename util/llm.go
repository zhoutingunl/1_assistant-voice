package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	myredis "github.com/Awaken1119/assistant-voice/cache/redis"
	"github.com/gomodule/redigo/redis"

	"github.com/Awaken1119/assistant-voice/appAction"
)

func llmServer(context string) (string, error) {
	// 若没有配置环境变量，可用百炼API Key将下行替换为：apiKey := "sk-xxx"。但不建议在生产环境中直接将API Key硬编码到代码中，以减少API Key泄露风险。
	apiKey := "sk-145f4423b1944ad181f020d7eacb95fc"
	appId := "ac3deb2cf3eb4b91b181407ee98c2298" // 替换为实际的应用 ID

	url := fmt.Sprintf("https://dashscope.aliyuncs.com/api/v1/apps/%s/completion", appId)

	// 创建请求体
	requestBody := map[string]interface{}{
		"input": map[string]string{
			"prompt": context,
		},
		"parameters": map[string]interface{}{},
		"debug":      map[string]interface{}{},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to marshal JSON: %v\n", err))
	}

	// 创建 HTTP POST 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to create request: %v\n", err))
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to send request: %v\n", err))
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to read response: %v\n", err))
	}

	llmResp := make(map[string]interface{})
	err = json.Unmarshal(body, &llmResp)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to unmarshal JSON: %v\n", err))
	}
	output, _ := llmResp["output"].(map[string]interface{})
	text, _ := output["text"].(string)
	// 处理响应
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("Request failed with status code: %d\n", resp.StatusCode))
	}
	return text, nil
}

func outputServer(output string) (string, error) {
	if strings.Contains(output, "<action>") {
		var actionType map[string]interface{}
		// 构建正则表达式模式
		pattern := `(?s)` + regexp.QuoteMeta("<action>") + "(.*?)" + regexp.QuoteMeta("</action>")
		re := regexp.MustCompile(pattern)
		action := re.FindAllStringSubmatch(output, -1)
		err := json.Unmarshal([]byte(action[0][1]), &actionType)
		acType := actionType["action"].(string)
		switch acType {
		case "open_app":
			output, err = appAction.OpenApp(action[0][1])
			if err != nil {
				return "", err
			}
			break
		case "write_text":
			output, err = appAction.WriteText(action[0][1])
			if err != nil {
				return "", err
			}
			break
		case "search_web":
			output, err = appAction.SearchWeb(action[0][1])
			if err != nil {
				return "", err
			}
			break
		case "open_file":
			output, err = appAction.OpenFile(action[0][1])
			if err != nil {
				return "", err
			}
			break
		case "play_music":
			output, err = appAction.PlayMusic(action[0][1])
			if err != nil {
				return "", err
			}
			break
		case "close_app":
			output, err = appAction.CloseApp(action[0][1])
			if err != nil {
				return "", err
			}
			break

		default:
			return "", errors.New(fmt.Sprintf("Unknown action type: %s\n", acType))
		}
		return output, nil
	}

	return "", errors.New(fmt.Sprintf("Output does not contain an answer. Output: %v", output))
}

func LlmRun(context string) (string, error) {
	rConn := myredis.InitRedis().Get()
	defer rConn.Close()
	appNameList, _ := redis.String(rConn.Do("GET", "appList"))
	for {
		question := fmt.Sprintf(context + "\n" + "<applist>" + appNameList + "</applist>")
		output, err := llmServer(question)
		if err != nil {
			return "", err
		}
		fmt.Println(output)
		if strings.Contains(output, "<final_answer>") {
			// 构建正则表达式模式
			pattern := `(?s)` + regexp.QuoteMeta("<final_answer>") + "(.*?)" + regexp.QuoteMeta("</final_answer>")
			re := regexp.MustCompile(pattern)
			answer := re.FindAllStringSubmatch(output, -1)
			return answer[0][1], nil
		}
		output, err = outputServer(output)
		if err != nil {
			return "", err
		}
		context = fmt.Sprintf("<observation>" + output + "</observation>")
	}
}
