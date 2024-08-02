package lib

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
)

// 使用者相關工具

// 產生唯一 UserID
func (tools *Tools) NewUserID() string {
	return uuid.New().String()
}

// 產生登入 SessionID
func (tools *Tools) NewSessionID(UserID string) (string, error) {
	// 生成隨機數
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)

	if err != nil {
		return "", err
	}
	// UserID + 當前時間戳 + 隨機數
	data := fmt.Sprintf("%s%d%s", UserID, tools.NowUnix(), randomBytes)
	// 使用 SHA256
	hash := sha256.Sum256([]byte(data))
	// 用 base64 加密, 得出長度約為43位數
	sessionID := base64.StdEncoding.EncodeToString(hash[:])

	return sessionID, nil
}

// 使用者密碼加密
func (tools *Tools) PwdEncode(pwd string) string {
	hash := sha256.Sum256([]byte(tools.config.PwdSalt + pwd))
	// 用 base64 加密，得出長度約為44位數
	return base64.StdEncoding.EncodeToString(hash[:])
}

// 產生帳號
func (tools *Tools) NewAccount(role int32) string {
	return "ACC" + tools.RandCharset(10)
}

// 產生邀請碼
func (tools *Tools) NewInvitationCode() string {
	// 邀請碼 10碼
	return tools.RandCharset(10)
}
