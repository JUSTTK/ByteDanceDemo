# 评论模块修复文档

## 文档信息

- **模块**: 评论模块 (Comment Module)
- **修复日期**: 2026-04-28
- **影响范围**: `cmd/api/service.go`, `controller/comment.go`, `service/commentServiceImpl.go`, `service/commentService.go`, `repository/commentDao.go`, `utils/sensetive/sensitiveFilter.go`
- **优先级**: P0 - 严重问题修复

---

## 问题概述

在对 ByteDanceDemo 项目的评论模块进行代码审查后，发现以下严重问题：

| 问题编号 | 问题描述 | 严重程度 | 影响 |
|---------|---------|---------|------|
| C-P0-1 | 敏感词过滤器重复初始化 | 高 | 每次请求都重新加载字典，性能问题 |
| C-P0-2 | GetCommentList数据竞争 | 高 | 多个goroutine并发append导致数据丢失/panic |
| C-P0-3 | 评论删除无权限检查 | 高 | 任何用户都可以删除任何人的评论 |

---

## 详细问题分析

### C-P0-1: 敏感词过滤器重复初始化

**位置**: `controller/comment.go:50`

**问题描述**:
每次发布评论请求都会调用 `sensetive.InitFilter()`，该函数会读取敏感词字典文件。

**问题代码** (已修复):
```go
if actionType == 1 {
    content := c.PostForm("comment_text")
    sensetive.InitFilter()  // 每次请求都调用！
    content = sensetive.Filter.Replace(content, '#')
```

**影响**:
- 每次请求都进行文件I/O操作
- 频繁重复初始化影响性能
- 类似点赞模块的问题

---

### C-P0-2: GetCommentList数据竞争

**位置**: `service/commentServiceImpl.go:156-169`

**问题描述**:
多个goroutine并发append到同一个slice，违反Go并发安全原则。

**问题代码** (已修复):
```go
commentInfoList := make([]Comment, 0, n)
var wg sync.WaitGroup
wg.Add(n)
for _, comment := range plainCommentList {
    var commentData Comment
    go func(comment model.Comment) {
        commentService.CombineComment(&commentData, &comment)
        commentInfoList = append(commentInfoList, commentData)  // 数据竞争！
        wg.Done()
    }(*comment)
}
```

**影响**:
- 可能导致数据丢失
- 可能触发panic
- 返回的数据不一致

---

### C-P0-3: 评论删除无权限检查

**位置**: `controller/comment.go:79-101`

**问题描述**:
删除评论时不检查用户是否是评论作者。

**问题代码** (已修复):
```go
} else {
    commentId, err := strconv.ParseInt(c.PostForm("comment_id"), 10, 64)
    // 没有检查用户是否是评论作者！
    err = commentService.DeleteCommentAction(commentId)
```

**影响**:
- 严重安全漏洞
- 任何用户都可以删除任何人的评论
- 数据被恶意破坏

---

## 实施详情

### Step 1: 修复敏感词过滤器单例模式

**文件**: `utils/sensetive/sensitiveFilter.go`

**修改前**:
```go
var Filter *sensitive.Filter

func InitFilter() {
    Filter = sensitive.New()
    err := Filter.LoadWordDict(WordDictPath)
    if err != nil {
        log.Println("InitFilter Fail,Err=" + err.Error())
    }
}
```

**修改后**:
```go
var (
    Filter     *sensitive.Filter
    filterOnce sync.Once
)

func InitFilter() {
    filterOnce.Do(func() {
        Filter = sensitive.New()
        err := Filter.LoadWordDict(WordDictPath)
        if err != nil {
            log.Println("InitFilter Fail,Err=" + err.Error())
        }
    })
}

func GetFilter() *sensitive.Filter {
    InitFilter()
    return Filter
}
```

---

### Step 2: 应用启动初始化

**文件**: `cmd/api/service.go`

**修改内容**: 在PreRun中添加敏感词过滤器初始化

```go
PreRun: func(cmd *cobra.Command, args []string) {
    config2.Init(config)
    log.InitLogger(mode)
    mysql.Init()
    redis.InitRedis()
    redis2.Init()
    redis3.Init()
    rabbitmq.InitRabbitMQ()
    rabbitmq.InitCommentRabbitMQ()
    rabbitmq.InitFollowRabbitMQ()
    dao.SetDefault(mysql.DB)
    service.StartFavoriteWorkerPool()
    sensetive.InitFilter()  // 添加敏感词过滤器初始化
},
```

---

### Step 3: 修改Controller

**文件**: `controller/comment.go`

**修改1**: 移除重复的InitFilter调用

```go
// 修改前
if actionType == 1 {
    content := c.PostForm("comment_text")
    sensetive.InitFilter()  // 删除这一行
    content = sensetive.Filter.Replace(content, '#')
    // ...
}

// 修改后
if actionType == 1 {
    content := c.PostForm("comment_text")
    content = sensetive.GetFilter().Replace(content, '#')  // 使用GetFilter()
    // ...
}
```

**修改2**: 添加权限检查

```go
// 修改前
err = commentService.DeleteCommentAction(commentId)

// 修改后
err = commentService.DeleteCommentAction(commentId, userId)
```

---

### Step 4: 创建CommentWorkerPool

**文件**: `service/commentServiceImpl.go`

**新增内容**: 在文件末尾添加Worker Pool

```go
// CommentWorkerPool 评论数据处理工作池
type CommentWorkerPool struct {
    workerCount int
}

var commentPoolInstance *CommentWorkerPool
var commentPoolOnce sync.Once

func GetCommentWorkerPool() *CommentWorkerPool {
    commentPoolOnce.Do(func() {
        commentPoolInstance = &CommentWorkerPool{
            workerCount: 10,
        }
    })
    return commentPoolInstance
}

// ProcessComments 并发处理评论列表（使用worker pool模式）
func (p *CommentWorkerPool) ProcessComments(
    plainComments []*model.Comment,
    videoId int64,
) ([]Comment, error) {
    n := len(plainComments)
    if n == 0 {
        return nil, nil
    }

    resultChan := make(chan Comment, n)
    var wg sync.WaitGroup
    wg.Add(n)

    sem := make(chan struct{}, p.workerCount)

    for _, comment := range plainComments {
        go func(c *model.Comment) {
            sem <- struct{}{}
            defer func() {
                <-sem
                wg.Done()
            }()

            var commentData Comment
            commentService := GetCommentServiceInstance()
            err := commentService.CombineComment(&commentData, c)
            if err != nil {
                log.Println("CombineComment error:", err)
                return
            }

            resultChan <- commentData

            videoIdToStr := strconv.FormatInt(videoId, 10)
            commentIdToStr := strconv.FormatInt(c.ID, 10)
            go insertRedisVCId(videoIdToStr, commentIdToStr, commentData)
        }(comment)
    }

    go func() {
        wg.Wait()
        close(resultChan)
    }()

    commentInfoList := make([]Comment, 0, n)
    for commentData := range resultChan {
        commentInfoList = append(commentInfoList, commentData)
    }

    sort.Sort(CommentSlice(commentInfoList))
    return commentInfoList, nil
}
```

---

### Step 5: 修改GetCommentList

**文件**: `service/commentServiceImpl.go`

**修改内容**: 使用Worker Pool替代原有的并发逻辑

```go
func (commentService *CommentServiceImpl) GetCommentList(videoId int64, userId int64) ([]Comment, error) {
    videoIdToStr := strconv.FormatInt(videoId, 10)
    cnt, err := redis.RdbVCid.SCard(redis.Ctx, videoIdToStr).Result()
    if err != nil {
        log.Println("SCard failed", err)
    }

    // 缓存命中
    if cnt > 0 {
        return getCommentsFromCache(videoIdToStr)
    }

    // 缓存未命中，查询数据库
    plainCommentList, err := repository.GetCommentList(videoId)
    if err != nil {
        log.Println(err.Error())
        return nil, err
    }

    // 使用Worker Pool处理评论列表
    pool := GetCommentWorkerPool()
    commentInfoList, err := pool.ProcessComments(plainCommentList, videoId)
    if err != nil {
        log.Println("ProcessComments error:", err)
        return nil, err
    }

    log.Println("get commentList success")
    return commentInfoList, nil
}
```

---

### Step 6: 修复权限检查

**文件**: `service/commentServiceImpl.go`

**修改内容**: 添加权限验证逻辑

```go
func (commentService *CommentServiceImpl) DeleteCommentAction(commentId int64, userId int64) error {
    // 先查询评论，检查用户权限
    plainComment, err := repository.GetCommentById(commentId)
    if err != nil {
        return fmt.Errorf("comment not found")
    }

    // 检查用户是否是评论作者
    if plainComment.UserID != userId {
        return fmt.Errorf("permission denied: you can only delete your own comments")
    }

    // 原有的删除逻辑...
}
```

---

### Step 7: 添加辅助方法

**文件**: `service/commentServiceImpl.go`

**新增**: `getCommentsFromCache` 方法

```go
// getCommentsFromCache 从缓存中获取评论列表
func getCommentsFromCache(videoIdToStr string) ([]Comment, error) {
    var commentInfoList []Comment
    commentIdStringList, err := redis.RdbVCid.SMembers(redis.Ctx, videoIdToStr).Result()
    if err != nil {
        log.Println("read redis vId failed", err)
        return nil, err
    }

    for _, commentIdString := range commentIdStringList {
        var commentData Comment
        commentString, err := redis.RdbCIdComment.Get(redis.Ctx, commentIdString).Result()
        if err != nil {
            log.Println("get comment from redis failed", err)
            continue
        }
        b := []byte(commentString)
        err = json.Unmarshal(b, &commentData)
        if err != nil {
            log.Println("unmarshal failed", err)
            continue
        }
        commentInfoList = append(commentInfoList, commentData)
    }

    log.Println("从redis读取的评论列表")
    sort.Sort(CommentSlice(commentInfoList))
    return commentInfoList, nil
}
```

---

### Step 8: 添加DAO方法

**文件**: `repository/commentDao.go`

**新增**: `GetCommentById` 方法

```go
func GetCommentById(commentId int64) (*model.Comment, error) {
    c := dao.Comment
    comment, err := c.Where(c.ID.Eq(commentId)).First()
    if err != nil {
        return nil, err
    }
    return comment, nil
}
```

---

### Step 9: 更新接口定义

**文件**: `service/commentService.go`

**修改内容**: 更新DeleteCommentAction方法签名

```go
type CommentService interface {
    GetCommentCnt(videoId int64) (int64, error)
    CommentAction(comment model.Comment) (Comment, error)
    DeleteCommentAction(commentId int64, userId int64) error  // 添加userId参数
    GetCommentList(videoId int64, userId int64) ([]Comment, error)
}
```

---

## 文件修改清单

| 文件 | 修改类型 | 说明 |
|------|---------|------|
| `utils/sensetive/sensitiveFilter.go` | 重构 | 改为单例模式 |
| `cmd/api/service.go` | 新增 | 添加敏感词过滤器初始化 |
| `controller/comment.go` | 修改 | 移除InitFilter调用，添加权限检查 |
| `service/commentServiceImpl.go` | 重构 | 添加Worker Pool，修改GetCommentList，添加权限检查 |
| `service/commentService.go` | 修改 | 更新DeleteCommentAction接口签名 |
| `repository/commentDao.go` | 新增 | 添加GetCommentById方法 |

---

## 向后兼容性保证

| 功能 | 兼容性 | 说明 |
|------|--------|------|
| API接口签名 | ✓ | `/comment/action/` 接口不变 |
| 请求参数 | ✓ | 格式不变 |
| 响应格式 | ✓ | 格式不变 |
| 缓存逻辑 | ✓ | Redis缓存保持不变 |
| 消息队列 | ✓ | RabbitMQ处理保持不变 |

**注意**: `DeleteCommentAction` 方法签名发生变化，但这只在内部使用，不影响外部API。

---

## 验证方法

### 单元测试

- [ ] 敏感词过滤正确工作
- [ ] 权限检查正确（用户只能删除自己的评论）
- [ ] Worker Pool正确处理评论列表
- [ ] 无数据竞争

### 集成测试

- [ ] 发布评论成功
- [ ] 删除自己评论成功
- [ ] 删除他人评论失败（返回权限错误）
- [ ] 评论列表正确返回
- [ ] Redis缓存正确更新

### 并发测试

- [ ] 多用户同时评论同一视频
- [ ] 多用户同时删除评论
- [ ] 高并发下GetCommentList无数据竞争
- [ ] 高并发下无goroutine泄漏

---

## 代码对比

### 并发处理对比

**修改前** (有数据竞争):
```go
commentInfoList := make([]Comment, 0, n)
var wg sync.WaitGroup
wg.Add(n)
for _, comment := range plainCommentList {
    var commentData Comment
    go func(comment model.Comment) {
        commentService.CombineComment(&commentData, &comment)
        commentInfoList = append(commentInfoList, commentData)  // 数据竞争！
        wg.Done()
    }(*comment)
}
wg.Wait()
```

**修改后** (线程安全):
```go
pool := GetCommentWorkerPool()
commentInfoList, err := pool.ProcessComments(plainCommentList, videoId)
```

### Worker Pool内部实现

```go
// 使用channel收集结果，避免数据竞争
resultChan := make(chan Comment, n)
var wg sync.WaitGroup
wg.Add(n)

// 使用限流的goroutine池
sem := make(chan struct{}, p.workerCount)

for _, comment := range plainComments {
    go func(c *model.Comment) {
        sem <- struct{}{}  // 获取信号量
        defer func() {
            <-sem          // 释放信号量
            wg.Done()
        }()

        resultChan <- commentData  // 通过channel安全返回
    }(comment)
}

// 收集结果
commentInfoList := make([]Comment, 0, n)
for commentData := range resultChan {
    commentInfoList = append(commentInfoList, commentData)
}
```

---

## 风险评估

- **低风险**: 敏感词过滤器改为单例
- **中风险**: GetCommentList逻辑重构，需要充分测试
- **高影响**: 权限检查是新增功能，需要确保不会影响现有用户

---

## 后续建议

### 短期优化

1. 移除随机的LikeCount/TeaseCount，改为真实数据
2. 修复表名不一致问题（`repository/commentDao.go:22`）
3. 统一错误处理风格
4. 添加参数校验（评论长度、内容非空等）

### 中期优化

1. 实现评论分页功能
2. 添加评论点赞功能
3. 添加回复评论功能
4. 添加评论举报功能

### 长期优化

1. 评论内容审核机制
2. 评论搜索功能
3. 评论数据分析
4. 反垃圾评论系统

---

## 总结

本次修复成功解决了评论模块的3个严重问题：

1. ✅ 修复了敏感词过滤器重复初始化问题
2. ✅ 修复了GetCommentList数据竞争问题（使用Worker Pool + Channel）
3. ✅ 修复了评论删除无权限检查问题

修复后的代码：
- 性能更好（避免重复加载敏感词字典）
- 更安全（添加权限检查）
- 更可靠（使用线程安全的并发处理）

---

**文档版本**: 1.0
**最后更新**: 2026-04-28
**维护者**: Development Team
