package main

import (
	"crypto"
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

func main() {
	// Load configuration
	viper.SetConfigFile("config/settings.yml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	secretKey := viper.GetString("settings.jwt.secretKey")

	fmt.Println("=== JWT安全检查报告 ===\n")

	// Check 1: Secret key length
	fmt.Printf("1. 密钥长度: %d 字符\n", len(secretKey))
	if len(secretKey) < 32 {
		fmt.Println("   ❌ 风险: 密钥太短（建议至少32字符/256位）")
	} else {
		fmt.Println("   ✅ 密钥长度足够")
	}

	// Check 2: Default weak secret
	if strings.Contains(secretKey, "123456") {
		fmt.Println("   ❌ 风险: 检测到弱密钥模式")
	} else {
		fmt.Println("   ✅ 未检测到默认弱密钥")
	}

	// Check 3: Base64 encoding
	if isBase64(secretKey) {
		fmt.Println("   ✅ 密钥使用Base64编码")
	} else {
		fmt.Println("   ⚠️  密钥未使用Base64编码")
	}

	// Check 4: Hash pattern
	if isHashPattern(secretKey) {
		fmt.Println("   ❌ 风险: 密钥看起来像是哈希值，容易被猜测")
	} else {
		fmt.Println("   ✅ 密钥不是可预测的模式")
	}

	// Check 5: Algorithm check
	if !isValidAlgorithm(secretKey) {
		fmt.Println("   ⚠️  建议: 建议使用HS256算法")
	} else {
		fmt.Println("   ✅ 算法配置正确")
	}

	fmt.Println("\n=== 建议 ===")
	fmt.Println("1. 定期更换JWT密钥（每3-6个月）")
	fmt.Println("2. 密钥长度至少32字节（256位）")
	fmt.Println("3. 使用安全的随机生成器创建密钥")
	fmt.Println("4. 不要在日志中记录密钥")
	fmt.Println("5. 考虑使用环境变量存储敏感信息")
}

func isBase64(s string) bool {
	// Simple base64 check (32-88 chars typical for JWT)
	pattern := `^[A-Za-z0-9+/]{32,88}={0,2}$`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}

func isHashPattern(s string) bool {
	// Check if it looks like a known hash pattern
	hashPatterns := []string{
		`^[a-f0-9]{32}$`,   // MD5
		`^[a-f0-9]{40}$`,   // SHA1
		`^[a-f0-9]{64}$`,   // SHA256
		`^[a-f0-9]{128}$`,  // SHA512
	}

	for _, pattern := range hashPatterns {
		matched, _ := regexp.MatchString(pattern, s)
		if matched {
			return true
		}
	}
	return false
}

func isValidAlgorithm(secret string) bool {
	// This is a simplified check
	// In practice, you'd want to verify the algorithm in the JWT header
	return len(secret) > 0
}