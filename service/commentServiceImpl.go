package service

import (
	"bytedancedemo/config"
	"bytedancedemo/middleware/rabbitmq"
	"bytedancedemo/middleware/redis"
	"bytedancedemo/model"
	"bytedancedemo/repository"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"
)

type CommentServiceImpl struct {
	UserService
}

var (
	commentServiceImpl *CommentServiceImpl
	commentServiceOnce sync.Once
)

func GetCommentServiceInstance() *CommentServiceImpl {
	commentServiceOnce.Do(func() {
		commentServiceImpl = &CommentServiceImpl{
			&UserServiceImpl{},
		}
	})
	return commentServiceImpl
}

func (commentService *CommentServiceImpl) CommentAction(comment model.Comment) (Comment, error) {
	csi := GetCommentServiceInstance()
	commentRes, err := repository.InsertComment(comment)
	if err != nil {
		return Comment{}, err
	}
	//GetUserLoginInfoById(id int64) (User, error)
	user, err := csi.GetUserDetailsById(comment.UserID, nil)
	if err != nil {
		log.Println(err.Error())
		return Comment{}, err
	}
	if user == nil {
		return Comment{}, fmt.Errorf("User not found") // Return an appropriate error
	}
	// 随机数生成种子
	rand.Seed(time.Now().Unix())
	commentData := Comment{
		Id:         commentRes.ID,
		User:       user,
		Content:    commentRes.Content,
		CreateDate: commentRes.CreatedAt.Format(config.GO_STARTER_TIME),
		LikeCount:  int64(rand.Intn(100)),
		TeaseCount: int64(rand.Intn(100)),
	}
	// redis操作：将发表的评论id存入redis
	go func() {
		insertRedisVCId(strconv.FormatInt(comment.VideoID, 10), strconv.FormatInt(commentRes.ID, 10), commentData)
		log.Println("commentAction save in redis")
	}()

	return commentData, nil
}

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

	// redis操作：先查redis，若有则更新redis. 若不在redis中，直接走数据库删除，返回客户端。
	commentIdToStr := strconv.FormatInt(commentId, 10)
	n, err := redis.RdbCVid.Exists(redis.Ctx, commentIdToStr).Result()
	if err != nil {
		log.Println(err)
	}
	// 删除评论的消息队列
	commentDelMQ := rabbitmq.SimpleCommentDelMQ
	// 缓存有此id
	if n > 0 {
		// 根据commentId查出对应的videoId
		vid, err := redis.RdbCVid.Get(redis.Ctx, commentIdToStr).Result()
		if err != nil {
			log.Println("redisCV not found", err)
		}
		// 删除缓存,CV直接把key删除，VC只需要移除下面的value(commentId)。
		n1, err := redis.RdbCVid.Del(redis.Ctx, commentIdToStr).Result()
		if err != nil {
			log.Println("redisCV delete failed", err)
		}
		n2, err := redis.RdbCIdComment.Del(redis.Ctx, commentIdToStr).Result()
		if err != nil {
			log.Println("redisCIdComment delete failed", err)
		}
		n3, err := redis.RdbVCid.SRem(redis.Ctx, vid, commentIdToStr).Result()
		if err != nil {
			log.Println("redisVc Remove failed", err)
		}
		log.Println("del comment in redis successfully:", n1, n2, n3)
		// 写数据库
		err = commentDelMQ.PublishSimple(commentIdToStr)
		return err
	}
	// 不在缓存中，直接走数据库，在消息队列中执行
	err = commentDelMQ.PublishSimple(commentIdToStr)
	return err
}

func (commentService *CommentServiceImpl) GetCommentList(videoId int64, userId int64) ([]Comment, error) {
	//redis操作：先查缓存是否命中，若命中，取缓存中的之；否则去读数据库并更新缓存
	videoIdToStr := strconv.FormatInt(videoId, 10)
	cnt, err := redis.RdbVCid.SCard(redis.Ctx, videoIdToStr).Result()
	if err != nil {
		log.Println("SCard failed", err)
	}
	// 缓存中存在评论列表
	if cnt > 0 {
		return getCommentsFromCache(videoIdToStr)
	}
	// 评论不在缓存中，评论既不在缓存也不在数据库中
	// 先根据videoId查评论id，再查用户信息
	plainCommentList, err := repository.GetCommentList(videoId)
	// 拿评论出错
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	n := len(plainCommentList)
	//fmt.Println("视频评论的数量：", n)
	// 如果没有评论, 即评论不在数据库也不在缓存
	if n == 0 {
		return nil, nil
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

func (commentService *CommentServiceImpl) CombineComment(comment *Comment, plainComment *model.Comment) error {
	commentServiceNew := GetCommentServiceInstance()
	user, err := commentServiceNew.GetUserDetailsById(plainComment.UserID, nil)
	if err == nil {
		comment.User = user
	}
	// 随机数生成种子
	rand.Seed(time.Now().Unix())
	comment.Id = plainComment.ID
	comment.Content = plainComment.Content
	comment.CreateDate = plainComment.CreatedAt.Format(config.GO_STARTER_TIME)
	comment.LikeCount = int64(rand.Intn(100))
	comment.TeaseCount = int64(rand.Intn(100))
	return nil
}

func (commentService *CommentServiceImpl) GetCommentCnt(videoId int64) (int64, error) {
	videoIdToStr := strconv.FormatInt(videoId, 10)
	cnt, err := redis.RdbVCid.SCard(redis.Ctx, videoIdToStr).Result()
	if err != nil {
		log.Println("SCard failed", err)
	}
	// 如果在缓存中直接返回
	if cnt > 0 {
		log.Println("从redis读取的评论数量")
		return cnt, nil
	}
	return repository.GetCommentCnt(videoId)
}

// redis中存储videId与commentId对应关系
func insertRedisVCId(videoId string, commentId string, comment Comment) {
	_, err := redis.RdbVCid.SAdd(redis.Ctx, videoId, commentId).Result()
	if err != nil {
		log.Println("redis save fail:vId-cId")
		redis.RdbVCid.Del(redis.Ctx, videoId)
		return
	}
	// 设置键的有效期，为数据不一致情况兜底
	redis.RdbVCid.Expire(redis.Ctx, videoId, redis.ExpireTime)
	// 设置键的有效期，为数据不一致情况兜底
	_, err = redis.RdbCVid.Set(redis.Ctx, commentId, videoId, redis.ExpireTime).Result()
	if err != nil {
		log.Println("redis save fail:cId-vId")
		return
	}
	b, err := json.Marshal(comment)
	if err != nil {
		log.Println("serialize failed in redis save", err)
	}
	// 设置键的有效期，为数据不一致情况兜底
	_, err = redis.RdbCIdComment.Set(redis.Ctx, commentId, string(b), redis.ExpireTime).Result()
	if err != nil {
		log.Println("redis save fail:cId-comment")
		return
	}
}

// CommentSlice Golang实现任意类型sort函数的流程
type CommentSlice []Comment

func (commentSlice CommentSlice) Len() int {
	return len(commentSlice)
}

func (commentSlice CommentSlice) Less(i, j int) bool {
	return commentSlice[i].CreateDate > commentSlice[j].CreateDate
}

func (commentSlice CommentSlice) Swap(i, j int) {
	commentSlice[i], commentSlice[j] = commentSlice[j], commentSlice[i]
}

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
