// Package service @Author: youngalone [2023/8/8]
package service

import (
	"bytedancedemo/dao"
	"bytedancedemo/middleware/redis"
	"bytedancedemo/model"
	"bytedancedemo/utils"
	"bytedancedemo/utils/encryption"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"sync"
	"time"
)

type UserServiceImpl struct {
	// 这里需要关注模块 点赞模块 视频模块的配合
	VideoService
	FollowService
}

var (
	userServiceImpl *UserServiceImpl
	once            sync.Once
)

func GetUserServiceInstance() *UserServiceImpl {
	once.Do(func() {
		userServiceImpl = &UserServiceImpl{
			VideoService:  &VideoServiceImp{},
			FollowService: &FollowServiceImp{},
		}
	})
	return userServiceImpl
}

// getUserCacheKey 生成用户详情缓存key
func getUserCacheKey(userId int64) string {
	return fmt.Sprintf("user_details:%d", userId)
}

// getUserFromCache 从Redis获取用户详情缓存
func getUserFromCache(userId int64) (*User, error) {
	key := getUserCacheKey(userId)
	data, err := redis.RdbUserDetails.Get(redis.Ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		zap.L().Error("解析用户缓存失败", zap.Int64("userId", userId), zap.Error(err))
		return nil, err
	}

	return &user, nil
}

// setUserToCache 将用户详情写入Redis缓存
func setUserToCache(user *User) {
	key := getUserCacheKey(user.Id)
	data, err := json.Marshal(user)
	if err != nil {
		zap.L().Error("序列化用户数据失败", zap.Int64("userId", user.Id), zap.Error(err))
		return
	}

	ttl := time.Duration((10 + rand.Int63n(20)) * int64(time.Minute))
	if err := redis.RdbUserDetails.Set(redis.Ctx, key, data, ttl).Err(); err != nil {
		zap.L().Error("写入用户缓存失败", zap.Int64("userId", user.Id), zap.Error(err))
	}
}

// invalidateUserCache 使指定用户的缓存失效
func invalidateUserCache(userId int64) {
	key := getUserCacheKey(userId)
	if err := redis.RdbUserDetails.Del(redis.Ctx, key).Err(); err != nil {
		zap.L().Error("删除用户缓存失败", zap.Int64("userId", userId), zap.Error(err))
	}
}

func (usi *UserServiceImpl) InsertUser(user *model.User) (res *model.User, isSuccess bool) {
	u := dao.User
	err := u.Create(user)
	if err != nil {
		zap.L().Error("新增用户失败", zap.String("err", err.Error()))
		return nil, false
	}
	resList, _ := u.Where(u.Name.Eq(user.Name), u.Password.Eq(user.Password)).Find()
	return resList[0], true
}

func (usi *UserServiceImpl) GetUserBasicByPassword(username string, hashedPassword string) (user *model.User, isExist bool) {
	u := dao.User
	resList, err := u.Where(u.Name.Eq(username)).Find()
	if err != nil {
		zap.L().Error("查询用户失败", zap.String("err", err.Error()))
		return nil, false
	}
	if len(resList) == 0 {
		zap.L().Warn("未查询到用户", zap.Error(errors.New("用户名或密码错误")))
		return nil, false
	}

	// Compare hashed password using bcrypt
	if err := encryption.ComparePassword(resList[0].Password, hashedPassword); err != nil {
		zap.L().Warn("密码不正确", zap.Error(err))
		return nil, false
	}

	return resList[0], true
}

// GetUserDetailsById 异步获取用户详情，使用Worker Pool + Redis缓存
// 流程：
// 1. 先从缓存读取，命中则直接返回
// 2. 缓存未命中，提交任务到Worker Pool（异步处理）
// 3. Worker处理完成后写入缓存，下次请求直接从缓存读取
func (usi *UserServiceImpl) GetUserDetailsById(id int64, curID *int64) (*User, error) {
	// 1. 先从缓存读取
	if user, err := getUserFromCache(id); err == nil && user != nil {
		// 缓存命中，更新关注状态（需要实时）
		if curID != nil {
			userService := GetUserServiceInstance()
			if isFollow, err := userService.CheckIsFollowing(id, *curID); err == nil {
				user.IsFollow = isFollow
			}
		}
		return user, nil
	}

	// 2. 缓存未命中，提交异步任务到Worker Pool
	pool := GetUserWorkerPool()
	task := UserQueryTask{
		UserID: id,
		CurID:  curID,
	}

	result, err := pool.Submit(task)
	if err != nil {
		return nil, err
	}

	// 3. Worker处理完成后，结果已自动写入缓存（在processQuery中）
	return result, nil
}

// GetUserName 在user表中根据id查询用户姓名
func (usi *UserServiceImpl) GetUserName(userId int64) (string, error) {
	u := dao.User

	var userName string
	err := u.
		Where(u.ID.Eq(userId)).
		Pluck(u.Name, &userName)
	if err != nil {
		zap.L().Error("查询用户名出错", zap.Error(err))
		return "", err
	}

	return userName, nil
}

// UserQueryTask 用户查询任务
type UserQueryTask struct {
	UserID     int64
	CurID      *int64
	ResultChan chan<- *User
}

// UserWorkerPool 用户查询工作池（自适应扩缩容）
type UserWorkerPool struct {
	taskChan    chan UserQueryTask
	minWorkers  int              // 最小worker数量
	maxWorkers  int              // 最大worker数量
	workerCount int              // 当前worker数量
	quitChan    chan struct{}    // 关闭信号
	wg          sync.WaitGroup
	once        sync.Once
	scaleMutex  sync.Mutex       // 保护worker数量调整
}

var userPoolInstance *UserWorkerPool

// InitUserWorkerPool 初始化用户查询工作池（自适应）
func InitUserWorkerPool(minWorkers, maxWorkers int) {
	userPoolInstance.once.Do(func() {
		userPoolInstance = &UserWorkerPool{
			taskChan:    make(chan UserQueryTask, 10000), // 扩大缓冲区到10000
			minWorkers:  minWorkers,
			maxWorkers:  maxWorkers,
			workerCount: minWorkers,
			quitChan:    make(chan struct{}),
		}
		userPoolInstance.start()
		go userPoolInstance.monitor() // 启动监控协程
	})
}

// GetUserWorkerPool 获取工作池单例
func GetUserWorkerPool() *UserWorkerPool {
	if userPoolInstance == nil {
		InitUserWorkerPool(10, 100) // 默认10-100个worker
	}
	return userPoolInstance
}

// start 启动初始worker
func (p *UserWorkerPool) start() {
	for i := 0; i < p.minWorkers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

// monitor 监控队列长度，自动扩缩容
func (p *UserWorkerPool) monitor() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			queueLen := len(p.taskChan)
			p.scaleWorkers(queueLen)
		case <-p.quitChan:
			return
		}
	}
}

// scaleWorkers 根据队列长度动态调整worker数量
func (p *UserWorkerPool) scaleWorkers(queueLen int) {
	p.scaleMutex.Lock()
	defer p.scaleMutex.Unlock()

	// 扩容：队列长度超过当前worker数量的一半
	if queueLen > p.workerCount/2 && p.workerCount < p.maxWorkers {
		newWorkers := minInt(p.workerCount*2, p.maxWorkers)
		addCount := newWorkers - p.workerCount
		zap.L().Info("扩容Worker Pool",
			zap.Int("from", p.workerCount),
			zap.Int("to", newWorkers),
			zap.Int("queueLen", queueLen))

		for i := 0; i < addCount; i++ {
			p.wg.Add(1)
			go p.worker()
		}
		p.workerCount = newWorkers
	}

	// 缩容：队列长度为0且worker数量超过最小值
	if queueLen == 0 && p.workerCount > p.minWorkers {
		targetWorkers := maxInt(p.workerCount/2, p.minWorkers)
		zap.L().Info("缩容Worker Pool",
			zap.Int("from", p.workerCount),
			zap.Int("to", targetWorkers),
			zap.Int("queueLen", queueLen))
		// 缩容通过让部分worker自然退出
		for i := 0; i < p.workerCount-targetWorkers; i++ {
			p.quitChan <- struct{}{}
		}
		p.workerCount = targetWorkers
	}
}

// worker 处理查询任务（异步）
func (p *UserWorkerPool) worker() {
	defer p.wg.Done()
	for {
		select {
		case task := <-p.taskChan:
			result := p.processQuery(task)
			task.ResultChan <- result
		case <-p.quitChan:
			return
		}
	}
}

// processQuery 处理单个查询任务，完成后写入缓存
func (p *UserWorkerPool) processQuery(task UserQueryTask) *User {
	userService := GetUserServiceInstance()
	u := dao.User
	resList, err := u.Where(u.ID.Eq(task.UserID)).Find()
	if err != nil || len(resList) == 0 {
		zap.L().Error("查询用户失败", zap.Int64("userId", task.UserID), zap.Error(err))
		return nil
	}

	user := &User{
		Id:              resList[0].ID,
		Name:            resList[0].Name,
		Avatar:          resList[0].Avatar,
		BackgroundImage: resList[0].BackgroundImage,
		Signature:       resList[0].Signature,
	}

	if workCnt, err := userService.GetVideoCountByAuthorID(task.UserID); err == nil {
		user.WorkCount = workCnt
	}
	if cnt, err := userService.GetFollowingCnt(task.UserID); err == nil {
		user.FollowCount = cnt
	}
	if cnt, err := userService.GetFollowerCnt(task.UserID); err == nil {
		user.FollowerCount = cnt
	}
	if likes, err := GetUserTotalReceivedLikes([]int64{task.UserID}); err == nil {
		user.TotalFavorited = likes[task.UserID]
	}
	if favorites, err := utils.GetUserFavorites([]int64{task.UserID}); err == nil {
		user.FavoriteCount = favorites[task.UserID]
	}
	if task.CurID != nil {
		if isFollow, err := userService.CheckIsFollowing(task.UserID, *task.CurID); err == nil {
			user.IsFollow = isFollow
		}
	}

	// 处理完成后写入缓存（异步，不阻塞）
	go setUserToCache(user)

	return user
}

// Submit 提交查询任务
func (p *UserWorkerPool) Submit(task UserQueryTask) (*User, error) {
	resultChan := make(chan *User, 1)
	task.ResultChan = resultChan

	select {
	case p.taskChan <- task:
	case <-p.quitChan:
		return nil, fmt.Errorf("worker pool已关闭")
	default:
		// 队列满了，尝试触发扩容
		go func() {
			p.scaleWorkers(len(p.taskChan) + 100)
		}()
		return nil, fmt.Errorf("系统繁忙，请稍后重试")
	}

	result := <-resultChan
	return result, nil
}

// minInt 返回两个整数的最小值
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// maxInt 返回两个整数的最大值
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
