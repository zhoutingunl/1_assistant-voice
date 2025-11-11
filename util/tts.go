package util

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type AudioOutput struct {
	URL string `json:"url"`
}

type Output struct {
	Audio        AudioOutput `json:"audio"`
	FinishReason string      `json:"finish_reason"`
}

type StreamEvent struct {
	Output Output `json:"output"`
}

// 定义请求结构体
type TTSRequest struct {
	Model string `json:"model"`
	Input struct {
		Text         string `json:"text"`
		Voice        string `json:"voice"`
		LanguageType string `json:"language_type"`
	} `json:"input"`
}

func VoiceTts(text string) (string, error) {
	apiKey := "sk-145f4423b1944ad181f020d7eacb95fc"

	// 构建请求体
	reqBody := TTSRequest{
		Model: "qwen3-tts-flash",
	}
	reqBody.Input.Text = text
	reqBody.Input.Voice = "Cherry"
	reqBody.Input.LanguageType = "Chinese"

	// 序列化为 JSON
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, _ := http.NewRequest("POST", "https://dashscope.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation", strings.NewReader(string(bodyBytes)))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-DashScope-SSE", "enable")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var audioURL string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "data:") {
			data := strings.TrimPrefix(line, "data:")
			data = strings.TrimSpace(data)
			if data == "[DONE]" || data == "" {
				continue
			}

			var event StreamEvent
			if err := json.Unmarshal([]byte(data), &event); err != nil {
				continue
			}

			if event.Output.Audio.URL != "" {
				audioURL = event.Output.Audio.URL
				fmt.Println("✅ 检测到最终音频 URL：", audioURL)
				break
			}
		}
	}

	if audioURL == "" {
		return "", errors.New("\"⚠️ 未检测到音频 URL，请检查 API 响应格式或打印原始行。\"")
	}

	//// 下载音频文件
	//resp2, err := http.Get(audioURL)
	//if err != nil {
	//	return "", err
	//}
	//defer resp2.Body.Close()
	//
	//outFile := "output.wav"
	//file, err := os.Create(outFile)
	//if err != nil {
	//	return "", err
	//}
	//defer file.Close()
	//
	//_, err = io.Copy(file, resp2.Body)
	//if err != nil {
	//	return "", err
	//}
	//
	//fmt.Printf("🎵 音频文件已保存：%s\n", outFile)
	return audioURL, nil
}
