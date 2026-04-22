package controller

import (
	"bytedancedemo/model"
	"bytedancedemo/service"
	"bytedancedemo/utils/encryption"
	"bytedancedemo/utils/token"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"regexp"
	"time"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User service.User `json:"user"`
}

// ValidateUsername validates username format
func ValidateUsername(username string) bool {
	// Username must be 3-20 characters, only letters and numbers
	re := regexp.MustCompile(`^[a-zA-Z0-9]{3,20}$`)
	return re.MatchString(username)
}

// ValidatePassword validates password format
func ValidatePassword(password string) bool {
	// Password must be at least 6 characters, contain both letters and numbers
	if len(password) < 6 {
		return false
	}

	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	return hasLetter && hasNumber
}

func Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Validate input
	if !ValidateUsername(username) {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户名格式不正确"},
		})
		return
	}

	if !ValidatePassword(password) {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "密码必须至少6位，包含字母和数字"},
		})
		return
	}

	passwordKey, err := encryption.EncryptPassword(password)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "密码加密失败"},
		})
		return
	}

	usi := service.GetUserServiceInstance()
	if _, isExist := usi.GetUserBasicByPassword(username, passwordKey); isExist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户已存在"},
		})
	} else {
		user, ok := usi.InsertUser(&model.User{Name: username, Password: passwordKey, Role: "common_user"})
		if !ok {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "注册失败"},
			})
		} else {
			tokenString, err := token.GenerateToken(
				[]byte(viper.GetString("settings.jwt.secretKey")),
				token.Claims{
					UserID:   user.ID,
					UserName: user.Name,
					Role:     user.Role,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour * viper.GetDuration("settings.jwt.expirationTime")).Unix(),
					},
				},
			)
			if err != nil {
				c.JSON(http.StatusOK, UserLoginResponse{
					Response: Response{StatusCode: 1, StatusMsg: "token令牌签发失败"},
				})
			} else {
				c.JSON(http.StatusOK, UserLoginResponse{
					Response: Response{StatusCode: 0},
					UserId:   user.ID,
					Token:    tokenString,
				})
			}
		}
	}
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	usi := service.GetUserServiceInstance()
	user, isExist := usi.GetUserBasicByPassword(username, password)
	if !isExist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "用户不存在",
			},
		})
	} else {
		tokenString, err := token.GenerateToken([]byte(viper.GetString("settings.jwt.secretKey")), token.Claims{
			UserID:   user.ID,
			UserName: user.Name,
			Role:     user.Role,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * viper.GetDuration("settings.jwt.expirationTime")).Unix(),
			},
		})
		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "token令牌签发失败",
				},
			})
		} else {
			zap.L().Debug("登录成功", zap.Int64("userID", user.ID), zap.String("username", user.Name), zap.String("role", user.Role))
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{
					StatusCode: 0,
				},
				UserId: user.ID,
				Token:  tokenString,
			})
		}
	}
}

func UserInfo(c *gin.Context) {
	userID := c.GetInt64("user_id")
	usi := service.GetUserServiceInstance()
	user, err := usi.GetUserDetailsById(userID, nil)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户不存在"},
		})
	} else {
		zap.L().Debug("查询用户详情成功", zap.Any("user", user))
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     *user,
		})
	}
}
