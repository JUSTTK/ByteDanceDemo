package controller

import (
	"bytedancedemo/dao"
	"log"
	"net/http"
	"strconv"

	"bytedancedemo/service"

	"github.com/gin-gonic/gin"
)

// RelationActionResp 关注和取消关注需要返回结构。
type RelationActionResp struct {
	Response
}

type UserListResponse struct {
	Response
	UserList []service.User `json:"user_list"`
}

type FriendUserListResponse struct {
	Response
	FriendUserList []service.FriendUser `json:"user_list"`
}

// checkUserExists 检查用户是否存在
func checkUserExists(userId int64) bool {
	u := dao.User
	count, err := u.Where(u.ID.Eq(userId)).Count()
	if err != nil {
		log.Printf("Check user exists failed: %v", err)
		return false
	}
	return count > 0
}

// RelationAction 关注/取关操作
func RelationAction(c *gin.Context) {
	userId := c.GetInt64("user_id")
	toUserId, err2 := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	actionType, err3 := strconv.ParseInt(c.Query("action_type"), 10, 64)

	// 传入参数格式有问题
	if nil != err2 || nil != err3 || actionType < 1 || actionType > 2 {
		c.JSON(http.StatusOK, RelationActionResp{
			Response{
				StatusCode: -1,
				StatusMsg:  "请求参数格式错误",
			},
		})
		return
	}

	if userId == toUserId {
		c.JSON(http.StatusOK, RelationActionResp{
			Response{
				StatusCode: -1,
				StatusMsg:  "不能关注或取关自己",
			},
		})
		return
	}

	if !checkUserExists(toUserId) {
		c.JSON(http.StatusOK, RelationActionResp{
			Response{
				StatusCode: -1,
				StatusMsg:  "目标用户不存在",
			},
		})
		return
	}

	// 正常处理
	fsi := service.NewFSIInstance()
	switch {
	// 关注
	case 1 == actionType:
		_, err := fsi.FollowAction(userId, toUserId)
		if err != nil {
			log.Printf("FollowAction failed: %v", err)
			c.JSON(http.StatusOK, RelationActionResp{
				Response{
					StatusCode: -1,
					StatusMsg:  "关注失败",
				},
			})
			return
		}
	// 取关
	case 2 == actionType:
		_, err := fsi.CancelFollowAction(userId, toUserId)
		if err != nil {
			log.Printf("CancelFollowAction failed: %v", err)
			c.JSON(http.StatusOK, RelationActionResp{
				Response{
					StatusCode: -1,
					StatusMsg:  "取关失败",
				},
			})
			return
		}
	default:
		c.JSON(http.StatusOK, RelationActionResp{
			Response{
				StatusCode: -1,
				StatusMsg:  "无效的操作类型",
			},
		})
		return
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "操作成功"})
}

// FollowList 获取关注列表
func FollowList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response{
				StatusCode: -1,
				StatusMsg:  "请求参数格式错误",
			},
			nil,
		})
		return
	}

	if !checkUserExists(userId) {
		c.JSON(http.StatusOK, UserListResponse{
			Response{
				StatusCode: -1,
				StatusMsg:  "用户不存在",
			},
			nil,
		})
		return
	}

	fsi := service.NewFSIInstance()
	followings, err1 := fsi.GetFollowings(userId)
	if err1 != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response{
				StatusCode: -1,
				StatusMsg:  "获取关注列表失败",
			},
			nil,
		})
		return
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response{
			StatusCode: 0,
			StatusMsg:  "获取关注列表成功",
		},
		followings,
	})
}

// FollowerList 获取粉丝列表
func FollowerList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response{
				StatusCode: -1,
				StatusMsg:  "请求参数格式错误",
			},
			nil,
		})
		return
	}

	if !checkUserExists(userId) {
		c.JSON(http.StatusOK, UserListResponse{
			Response{
				StatusCode: -1,
				StatusMsg:  "用户不存在",
			},
			nil,
		})
		return
	}

	fsi := service.NewFSIInstance()
	followers, err1 := fsi.GetFollowers(userId)
	if err1 != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response{
				StatusCode: -1,
				StatusMsg:  "获取粉丝列表失败",
			},
			nil,
		})
		return
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response{
			StatusCode: 0,
			StatusMsg:  "获取粉丝列表成功",
		},
		followers,
	})
}

// FriendList 获取好友列表
func FriendList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusOK, FriendUserListResponse{
			Response{
				StatusCode: -1,
				StatusMsg:  "请求参数格式错误",
			},
			nil,
		})
		return
	}

	if !checkUserExists(userId) {
		c.JSON(http.StatusOK, FriendUserListResponse{
			Response{
				StatusCode: -1,
				StatusMsg:  "用户不存在",
			},
			nil,
		})
		return
	}

	fsi := service.NewFSIInstance()
	friends, err1 := fsi.GetFriends(userId)
	if err1 != nil {
		c.JSON(http.StatusOK, FriendUserListResponse{
			Response{
				StatusCode: -1,
				StatusMsg:  "获取好友列表失败",
			},
			nil,
		})
		return
	}

	c.JSON(http.StatusOK, FriendUserListResponse{
		Response{
			StatusCode: 0,
			StatusMsg:  "获取好友列表成功",
		},
		friends,
	})
}
