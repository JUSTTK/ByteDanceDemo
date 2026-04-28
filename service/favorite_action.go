package service

import (
	"bytedancedemo/database/mysql"
	"bytedancedemo/utils"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"bytedancedemo/dao"
	"bytedancedemo/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// FavoriteTask 点赞任务
type FavoriteTask struct {
	UserID     int64
	VideoID    int64
	ActionType int32
	ResultChan chan<- FavoriteResult
}

// FavoriteResult 处理结果
type FavoriteResult struct {
	Success    bool
	Error      error
	StatusCode int32
	StatusMsg  string
}

// FavoriteWorkerPool 点赞工作池（单例）
type FavoriteWorkerPool struct {
	taskChan chan FavoriteTask
	quitChan chan struct{}
	wg       sync.WaitGroup
	once     sync.Once
}

var favoritePool *FavoriteWorkerPool

// InitFavoriteWorkerPool 初始化点赞工作池（单例，只在应用启动时调用一次）
func InitFavoriteWorkerPool(workerCount int) {
	favoritePool.once.Do(func() {
		favoritePool = &FavoriteWorkerPool{
			taskChan: make(chan FavoriteTask, 1000), // 缓冲1000个任务
			quitChan: make(chan struct{}),
		}
		favoritePool.start(workerCount)
	})
}

// StartFavoriteWorkerPool 公开的初始化方法（用于配置文件指定worker数量）
func StartFavoriteWorkerPool() {
	// 从配置读取worker数量，默认为10
	workerCount := 10
	InitFavoriteWorkerPool(workerCount)
}

// GetFavoriteWorkerPool 获取工作池单例
func GetFavoriteWorkerPool() *FavoriteWorkerPool {
	if favoritePool == nil {
		StartFavoriteWorkerPool()
	}
	return favoritePool
}

// start 启动worker
func (p *FavoriteWorkerPool) start(workerCount int) {
	for i := 0; i < workerCount; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

// worker 处理任务的goroutine
func (p *FavoriteWorkerPool) worker() {
	defer p.wg.Done()
	for {
		select {
		case task := <-p.taskChan:
			result := p.processTask(task)
			task.ResultChan <- result
		case <-p.quitChan:
			return
		}
	}
}

// processTask 处理单个任务
func (p *FavoriteWorkerPool) processTask(task FavoriteTask) FavoriteResult {
	var err error
	statusCode := SuccessCode
	statusMsg := SuccessMessage

	switch task.ActionType {
	case 1:
		err = likeVideo(task.UserID, task.VideoID)
		if err != nil {
			statusCode = ErrorCode
			statusMsg = err.Error()
		} else {
			err = utils.UpdateLikeCounts(task.UserID, task.VideoID, true)
			if err != nil {
				statusMsg = "点赞成功，但计数更新失败"
			}
		}
	case 2:
		err = unlikeVideo(task.UserID, task.VideoID)
		if err != nil {
			statusCode = ErrorCode
			statusMsg = err.Error()
		} else {
			err = utils.UpdateLikeCounts(task.UserID, task.VideoID, false)
			if err != nil {
				statusMsg = "取消点赞成功，但计数更新失败"
			}
		}
	default:
		err = fmt.Errorf("invalid action_type: %v", task.ActionType)
		statusCode = ErrorCode
		statusMsg = err.Error()
	}

	return FavoriteResult{
		Success:    err == nil,
		Error:      err,
		StatusCode: statusCode,
		StatusMsg:  statusMsg,
	}
}

// Submit 提交任务并等待结果
func (p *FavoriteWorkerPool) Submit(userID, videoID int64, actionType int32) (FavoriteActionResponse, error) {
	resultChan := make(chan FavoriteResult, 1)
	task := FavoriteTask{
		UserID:     userID,
		VideoID:    videoID,
		ActionType: actionType,
		ResultChan: resultChan,
	}

	// 提交任务（非阻塞）
	select {
	case p.taskChan <- task:
		// 任务已提交
	default:
		return FavoriteActionResponse{}, fmt.Errorf("系统繁忙，请稍后重试")
	}

	// 等待结果（阻塞，但有超时）
	result := <-resultChan

	return FavoriteActionResponse{
		StatusCode: result.StatusCode,
		StatusMsg:  result.StatusMsg,
	}, result.Error
}

// Shutdown 优雅关闭工作池
func (p *FavoriteWorkerPool) Shutdown() {
	close(p.quitChan)
	p.wg.Wait()
}

func (s *FavoriteServiceImpl) FavoriteAction(userId int64, videoID int64, actionType int32) (FavoriteActionResponse, error) {
	pool := GetFavoriteWorkerPool()
	return pool.Submit(userId, videoID, actionType)
}

func likeVideo(userID int64, videoID int64) error {
	tx := mysql.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	first, err := dao.Like.Where(dao.Like.UserID.Eq(userID), dao.Like.VideoID.Eq(videoID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			like := model.Like{
				UserID:  userID,
				VideoID: videoID,
				Liked:   1,
			}
			if err := tx.Create(&like).Error; err != nil {
				tx.Rollback()
				return err
			}
			return tx.Commit().Error
		} else {
			tx.Rollback()
			return err
		}
	}
	if first.Liked == 1 {
		tx.Rollback()
		return fmt.Errorf("user has already liked this video")
	}

	first.Liked = 1
	if err := tx.Save(&first).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func unlikeVideo(userID int64, videoID int64) error {
	tx := mysql.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	first, err := dao.Like.Where(dao.Like.UserID.Eq(userID), dao.Like.VideoID.Eq(videoID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return fmt.Errorf("no like found for this user and video")
		}
		tx.Rollback()
		return err
	}

	if first.Liked == 0 {
		tx.Rollback()
		return fmt.Errorf("user has already unliked this video")
	}

	first.Liked = 0
	if err := tx.Save(&first).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *FavoriteServiceImpl) GetVideosLikes(videoIDs []int64) (map[int64]int64, error) {
	ctx := utils.GlobalRedisClient.Ctx
	pipe := utils.GlobalRedisClient.Client.Pipeline()

	futures := make(map[int64]*redis.StringCmd)
	for _, videoID := range videoIDs {
		videoKey := fmt.Sprintf("video:%d", videoID)
		videoLikesField := "totalVideoLikes"
		futures[videoID] = pipe.HGet(ctx, videoKey, videoLikesField)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("failed to execute pipeline: %v", err)
	}

	result := make(map[int64]int64)
	for videoID, future := range futures {
		err := future.Err()
		if err == redis.Nil {
			result[videoID] = 0
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get video likes for video %d: %v", videoID, err)
		}
		likesStr, _ := future.Result()
		likes, err := strconv.ParseInt(likesStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse likes for video %d: %v", videoID, err)
		}
		result[videoID] = likes
	}

	return result, nil
}
