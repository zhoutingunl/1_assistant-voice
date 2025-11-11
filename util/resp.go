package util

type LoginResp struct {
	Success      bool   `json:"success"`      // 登录是否成功
	Message      string `json:"message"`      // 提示信息，失败时返回错误原因
	Token        string `json:"token"`        // 访问令牌
	RefreshToken string `json:"refreshToken"` // 刷新令牌
}

type Resp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Message string `json:"message"`
	}
}
