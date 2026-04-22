# ByteDanceDemo API 使用示例

本文档提供了 ByteDanceDemo API 的详细使用示例，包括 cURL、JavaScript、Python 和 Go 完整的 API 调用代码。

## 目录

- [用户认证](#用户认证)
- [视频管理](#视频管理)
- [社交功能](#社交功能)
- [消息系统](#消息系统)
- [通用说明](#通用说明)

## 用户认证

### 1. 用户注册

#### cURL

```bash
curl -X POST "http://localhost:8080/douyin/user/register/" \
  -F "username=testuser" \
  -F "password=secure123"
```

**响应示例**：
```json
{
  "status_code": 0,
  "status_msg": "注册成功",
  "user_id": 1
}
```

#### JavaScript

```javascript
fetch('http://localhost:8080/douyin/user/register/', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/x-www-form-urlencoded'
  },
  body: new URLSearchParams({
    username: 'testuser',
    password: 'secure123'
  })
})
.then(response => response.json())
.then(data => {
  console.log('注册成功:', data);
  console.log('用户 ID:', data.user_id);
  console.log('Token:', data.token);
})
.catch(error => {
  console.error('注册失败:', error);
});
```

#### Python

```python
import requests
import json

def register_user(username, password):
    url = "http://localhost:8080/douyin/user/register/"
    data = {
        "username": username,
        "password": password
    }
    
    response = requests.post(url, data=data)
    result = response.json()
    
    if result["status_code"] == 0:
        print("注册成功!")
        print(f"用户 ID: {result['user_id']}")
        print(f"Token: {result['token']}")
    else:
        print(f"注册失败: {result['status_msg']}")

# 使用示例
register_user("testuser", "secure123")
```

#### Go

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
)

func registerUser(username, password string) error {
    data := map[string]string{
        "username": username,
        "password": password,
    }
    jsonData, err := json.Marshal(data)
    if err != nil {
        return err
    }
    
    resp, err := http.Post(
        "http://localhost:8080/douyin/user/register/",
        "application/x-www-form-urlencoded",
        bytes.NewBuffer(jsonData),
    )
    
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    
    fmt.Printf("状态码: %v\n", result["status_code"])
    fmt.Printf("状态消息: %v\n", result["status_msg"])
    fmt.Printf("用户 ID: %v\n", result["user_id"])
    fmt.Printf("Token: %v\n", result["token"])
    
    return nil
}

func main() {
    err := registerUser("testuser", "secure123")
    if err != nil {
        fmt.Printf("错误: %v\n", err)
    }
}
```

### 2. 用户登录

#### cURL

```bash
curl -X POST "http://localhost:8080/douyin/user/login/" \
  -F "username=testuser" \
  -F "password=secure123"
```

**响应示例**：
```json
{
  "status_code": 0,
  "status_msg": "登录成功",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ..."
}
```

#### JavaScript

```javascript
// 使用获取的 token
const token = localStorage.getItem('token') || 'YOUR_JWT_TOKEN';

fetch('http://localhost:8080/douyin/user/', {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${token}`
  }
})
.then(response => response.json())
.then(data => {
  console.log('用户信息:', data.user);
})
.catch(error => {
  console.error('获取用户信息失败:', error);
});
```

#### Python

```python
import requests

def get_user_info(token):
    url = "http://localhost:8080/douyin/user/"
    headers = {
        "Authorization": f"Bearer {token}"
    }
    
    response = requests.get(url, headers=headers)
    result = response.json()
    
    if result["status_code"] == 0:
        print("用户信息:")
        print(f"  ID: {result['user']['id']}")
        print(f" 名字: {result['user']['name']}")
        print(f" 头像: {result['user']['avatar']}")
        print(f" 签名数: {result['user']['follow_count']}")
        print(f" 粉丝数: {result['user']['follower_count']}")
        print(f" 获赞数: {result['user']['favorite_count']}")
    else:
        print(f"失败: {result['status_msg']}")

# 使用示例
get_user_info("YOUR_JWT_TOKEN")
```

#### Go

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
)

type UserResponse struct {
    StatusCode int    `json:"status_code"`
    StatusMsg  string `json:"status_msg"`
    User      User   `json:"user"`
}

type User struct {
    ID          int64  `json:"id"`
    Name        string `json:"name"`
    Avatar      string `json:"avatar"`
    FollowCount int64  `json:"follow_count"`
    FollowerCount int64 `json:"follower_count"`
    FavoriteCount int64  `json:"favorite_count"`
}

func getUserInfo(token string) error {
    req, err := http.NewRequest("GET", "http://localhost:8080/douyin/user/", nil)
    if err != nil {
        return err
    }
    
    req.Header.Set("Authorization", "Bearer "+token)
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    var result UserResponse
    json.NewDecoder(resp.Body).Decode(&result)
    
    fmt.Printf("用户信息:\n")
    fmt.Printf("  ID: %d\n", result.User.ID)
    fmt.Printf(" 名字: %s\n", result.User.Name)
    fmt.Printf(" 头像: %s\n", result.User.Avatar)
    fmt.Printf(" 关注数: %d\n", result.User.FollowCount)
    fmt.Printf(" 粉丝数: %d\n", result.User.FollowerCount)
    fmt.Printf(" 获赞数: %d\n", result.User.FavoriteCount)
    
    return nil
}

func main() {
    token := "YOUR_JWT_TOKEN"
    err := getUserInfo(token)
    if err != nil {
        fmt.Printf("错误: %v\n", err)
    }
}
```

## 视频管理

### 1. 获取视频流

#### cURL

```bash
curl -X GET "http://localhost:8080/douyin/feed/?token=YOUR_TOKEN&latest_time=0"
```

**响应示例**：
```json
{
  "status_code": 0,
  "status_msg": "",
  "video_list": [
    {
      "id": 1,
      "author": {
        "id": 1,
        "name": "user1",
        "avatar": "http://example.com/avatar1.jpg",
        "follow_count": 10,
        "follower_count": 20
        "is_following": true
      },
      "play_url": "http://localhost:8080/static/video1.mp4",
      "cover_url": "http://localhost:8080/static/cover1.jpg",
      "favorite_count": 5,
      "comment_count": 3,
      "title": "我的第一个视频",
      "create_time": "2023-01-01T00:00:00Z"
    }
  ]
}
```

#### JavaScript

```javascript
const token = localStorage.getItem('token') || 'YOUR_JWT_TOKEN';

async function getVideoFeed() {
    const response = await fetch('http://localhost:8080/douyin/feed/', {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    
    const data = await response.json();
    
    console.log('获取到', data.video_list.length, '个视频');
    
    data.video_list.forEach(video => {
        console.log(`视频 ${video.id}: ${video.title}`);
        console.log(`作者: ${video.author.name}`);
        console.log(`点赞数: ${video.favorite_count}`);
        console.log(`评论数: ${video.comment_count}`);
    });
}
}
```

### 2. 发布视频

#### cURL

```bash
curl -X POST "http://localhost:8080/douyin/publish/action/" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "data=@/path/to/video.mp4" \
  -F "token=YOUR_JWT_TOKEN"
```

#### JavaScript - 支持进度上传

```javascript
const token = localStorage.getItem('token') || 'YOUR_JWT_TOKEN';

async function uploadVideo(file) {
    const formData = new FormData();
    formData.append('data', file);
    formData.append('token', token);
    
    const response = await fetch('http://localhost:8080/douyin/publish/action/', {
        method: 'POST',
        body: formData,
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    
    const data = await response.json();
    
    if (data.status_code === 0) {
        console.log('视频发布成功！');
        console.log('视频 URL:', data.video.play_url);
    } else {
        console.error('发布失败:', data.status_msg);
    }
}
```

## 社交功能

### 1. 点赞视频

#### cURL

```bash
curl -X POST "http://localhost:8080/douyin/favorite/action/" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "video_id=1" \
  -F "action_type=1" \
  -F "token=YOUR_JWT_TOKEN"
```

**action_type 说明**：
- 1: 点赞
- 2: 取消点赞

#### JavaScript

```javascript
const token = localStorage.getItem('token') || 'YOUR_JWT_TOKEN';

async function likeVideo(videoId) {
    const response = await fetch('http://localhost:8080/douyin/favorite/action/', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
            video_id: videoId,
            action_type: 1,  // 1 表示点赞
            token: token
        })
    });
    
    const data = await response.json();
    console.log('操作结果:', data.status_msg);
}
```

### 2. 获取点赞列表

#### cURL

```bash
curl -X GET "http://localhost:8080/douyin/favorite/list/?token=YOUR_JWT_TOKEN&user_id=1"
```

### 3. 检查视频点赞状态

通常在获取视频流时，响应中会包含 `is_favorite` 字段，可以直接使用。

## 评论系统

### 1. 发布评论

#### cURL

```bash
curl -X POST "http://localhost:8080/douyin/comment/action/" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "video_id=1" \
  -F "comment_text=这是一个很棒的视频！" \
  -F "token=YOUR_JWT_TOKEN"
```

#### JavaScript

```javascript
const token = localStorage.getItem('token') || 'YOUR_JWT_TOKEN';

async function postComment(videoId, commentText) {
    const response = await fetch('http://localhost:8080/douyin/comment/action/', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
            video_id: videoId,
            comment_text: commentText,
            token: token
        })
    });
    
    const data = await response.json();
    
    if (data.status_code === 0) {
        console.log('评论发布成功！');
        console.log('评论 ID:', data.comment.id);
    } else {
        console.error('评论失败:', data.status_msg);
    }
}
```

### 2. 获取视频评论

#### cURL

```bash
curl -X GET "http://localhost:8080/douyin/comment/list/?video_id=1&token=YOUR_JWT_TOKEN"
```

**响应示例**：
```json
{
  "status_code": 0,
  "status_msg": "",
  "comment_list": [
    {
      "id": 1,
      "user": {
        "id": 1,
        "name": "user1",
        "avatar": "http://example.com/avatar1.jpg"
      },
      "content": "这是一个很棒的视频！",
      "create_date": "2023-01-01T12:00:00Z"
    }
  ]
}
```

## 关注系统

### 1. 关注用户

#### cURL

```bash
curl -X POST "http://localhost:8080/douyin/relation/action/" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "to_user_id=2" \
  -F "action_type=1" \
  - F "token=YOUR_JWT_TOKEN"
```

**action_type 说明**：
- 1: 关注
- 2: 取消关注

#### JavaScript

```javascript
const token = localStorage.getItem('token') || 'YOUR_JWT_TOKEN';

async function followUser(userId) {
    const response = await fetch('http://localhost:8080/douyin/relation/action/', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
            to_user_id: userId,
            action_type: 1,  // 1 表示关注
            token: token
        })
    });
    
    const data = await response.json();
    console.log('操作结果:', data.status_msg);
}
```

### 2. 获取关注列表

#### cURL

```bash
curl -X GET "http://localhost:8080/douyin/relation/follow/list/?token=YOUR_JWT_TOKEN&user_id=1"
```

### 3. 获取粉丝列表

#### cURL

```bash
curl -X GET "http://localhost:8080/douyin/relation/follower/list/?token=YOUR_JWT_TOKEN&user_id=1"
```

### 4. 获取好友列表

#### cURL

```bash
curl -X GET "http://localhost:8080/douyin/relation/friend/list/?token=YOUR_JWT_TOKEN&user_id=1"
```

## 消息系统

### 1. 发送消息

#### cURL

```bash
curl -X POST "http://localhost:8080/douyin/message/action/" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "to_user_id=2" \
  -F "content=你好，这是一条消息！" \
  -F "token=YOUR_JWT_TOKEN"
```

#### JavaScript

```javascript
const token = localStorage.getItem('token') || 'YOUR_JWT_TOKEN';

async function sendMessage(toUserId, content) {
    const response = await fetch('http://localhost:8080/douyin/message/action/', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
            to_user_id: toUserId,
            content: content,
            token: token
        })
    });
    
    const data = await response.json();
    console.log('消息发送结果:', data.status_msg);
}
```

### 2. 获取聊天记录

#### cURL

```bash
curl -X GET "http://localhost:8080/douyin/message/chat/?to_user_id=2&token=YOUR_JWT_TOKEN"
```

**响应示例**：
```json
{
  "status_code": 0,
  "status_msg": "",
  "message_list": [
    {
      "id": 1,
      "from_user_id": 1,
      "to_user_id": 2,
      "content": "你好，这是一条消息！",
      "create_time": "2023-01-01T12:00:00Z"
    },
    {
      "id": 2,
      "from_user_id": 2,
      "to_user_id": 1,
      "content": "谢谢！",
      "create_time": "2023-01-01T12:05:00Z"
    }
  ]
}
```

## 通用说明

### 错误响应格式

所有 API 在失败时都会返回以下格式：

```json
{
  "status_code": -1,
  "status_msg": "错误描述"
}
```

常见错误码：
- `-1`: 未授权
- `-1`: 参数错误
- `-1`: 数据库错误

### 速率限制

默认限制：每分钟 50 个请求

超过限制时会返回：
```json
{
  "status_code": -1,
  "status_msg": "请求过于频繁"
}
```

### 认证失败

当 token 无效或过期时：
```json
{
  "status_code": -1,
  "status_msg": "认证失败"
}
```

### 分页

大多数列表接口支持分页参数：

```bash
# 第一页，每页 10 个
curl "http://localhost:8080/douyin/feed/?token=TOKEN&page=1&page_size=10"

# 第二页
curl "http://localhost:8080/douyin/feed/?token=TOKEN&page=2&page_size=10"
```

### 文件上传

文件上传使用 multipart/form-data 格式，最大文件大小限制为 50MB。

## 调试技巧

### 使用 curl 的 -v 选项

```bash
# 查看完整的请求和响应
curl -v "http://localhost:8080/douyin/feed/?token=TOKEN"

# 只看响应头
curl -I "http://localhost:8080/douyin/feed/?token=TOKEN"

# 保存响应到文件
curl -o response.json "http://localhost:8080/douyin/feed/?token=TOKEN"
```

### 使用 JSON 格式化输出

```bash
# 格式化输出 JSON
curl "http://localhost:8080/douyin/feed/?token=TOKEN" | jq .

# 美化输出（美化 JSON）
curl "http://localhost:8080/douyin/feed/?token=TOKEN" | jq '.'
```

### 保存 token

登录成功后，保存 token 以供后续请求使用：

```bash
# 保存到环境变量
export BYTEDANCEDEMO_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 保存到文件
echo "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." > token.txt
```

### 批量操作

使用 `jq` 或 `gojq` 处理 JSON 数据：

```bash
# 提取所有用户 ID
curl "http://localhost:8080/douyin/feed/?token=TOKEN" | jq '.video_list[].author.id'

# 计算总数
curl "http://localhost:8080/douyin/feed/?token=TOKEN" | jq '.video_list | length'
```

---

**最后更新**: 2026-04-21
