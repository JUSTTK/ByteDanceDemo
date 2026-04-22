package service

import (
	"bytedancedemo/dao"
	"bytedancedemo/middleware/rabbitmq"
	"bytedancedemo/middleware/redis"
	"bytedancedemo/model"
	"bytedancedemo/repository"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// FollowServiceImp 该结构体继承FollowService，MessageService，UserService接口。
type FollowServiceImp struct {
	MessageService
	FollowService
	UserService
}

var (
	followServiceImp  *FollowServiceImp //controller层通过该实例变量调用service的所有业务方法。
	followServiceOnce sync.Once         //限定该service对象为单例，节约内存。
)

func CacheTimeGenerator() time.Duration {
	// 先设置随机数 - 这里比较重要
	rand.Seed(time.Now().Unix())
	// 再设置缓存时间
	// 10 + [0~20) 分钟的随机时间
	return time.Duration((10 + rand.Int63n(20)) * int64(time.Minute))
}

// 将一个字符串数组 strArr 转换为对应的 int64 类型的整数数组
func convertToInt64Array(strArr []string) ([]int64, error) {
	int64Arr := make([]int64, len(strArr))
	for i, str := range strArr {
		int64Val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
		int64Arr[i] = int64Val
	}
	return int64Arr, nil
}

// NewFSIInstance 生成并返回FollowServiceImp结构体单例变量。
func NewFSIInstance() *FollowServiceImp {
	//Do 方法，该方法接收一个函数作为参数，并且保证该函数在程序运行过程中只会被执行一次。
	followServiceOnce.Do(
		func() {
			followServiceImp = &FollowServiceImp{
				UserService: &UserServiceImpl{},
			}
		})
	return followServiceImp
}

//-------------------------------------API IMPLEMENT--------------------------------------------

/*
	关注业务
*/

// FollowAction 关注操作的业务
func (followService *FollowServiceImp) FollowAction(userId int64, targetId int64) (bool, error) {

	followDao := repository.NewFollowDaoInstance()
	follow, err := followDao.FindEverFollowing(userId, targetId)
	if nil != err {
		return false, err
	}

	// 获取关注的消息队列
	followAddMQ := rabbitmq.SimpleFollowAddMQ

	var dbErr error
	// 曾经关注过，只需要update一下followed即可。
	if nil != follow {
		_, dbErr = followDao.UpdateFollowRelation(userId, targetId, int8(1))
		if dbErr != nil {
			log.Printf("UpdateFollowRelation failed: %v", dbErr)
			return false, dbErr
		}
	} else {
		_, dbErr = followDao.InsertFollowRelation(userId, targetId)
		if dbErr != nil {
			log.Printf("InsertFollowRelation failed: %v", dbErr)
			return false, dbErr
		}
	}

	// 数据库操作成功后，发送消息队列通知其他模块
	if nil != follow {
		err := followAddMQ.PublishSimpleFollow(fmt.Sprintf("%d-%d-%s", userId, targetId, "update"))
		if err != nil {
			log.Printf("Publish follow update message failed: %v", err)
		}
	} else {
		err = followAddMQ.PublishSimpleFollow(fmt.Sprintf("%d-%d-%s", userId, targetId, "insert"))
		if err != nil {
			log.Printf("Publish follow insert message failed: %v", err)
		}
	}

	// 更新缓存
	followService.AddToRDBWhenFollow(userId, targetId)
	return true, nil

}

func (followService *FollowServiceImp) AddToRDBWhenFollow(userId int64, targetId int64) {
	followDao := repository.NewFollowDaoInstance()
	// 尝试给following数据库追加user关注target的记录
	keyCnt1, err1 := redis.UserFollowings.Exists(redis.Ctx, strconv.FormatInt(userId, 10)).Result()

	if err1 != nil {
		log.Printf("Check Redis existence failed: %v", err1)
	}

	// 只判定键是否不存在，若不存在即从数据库导入
	if keyCnt1 <= 0 {
		userFollowingsId, _, err := followDao.GetFollowingsInfo(userId)
		if err != nil {
			log.Printf("GetFollowingsInfo failed: %v", err)
			return
		}
		ImportToRDBFollowing(userId, userFollowingsId)
	}
	// 数据库导入到redis结束后追加记录
	redis.UserFollowings.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), targetId)

	// 尝试给follower数据库追加target的粉丝有user的记录
	keyCnt2, err2 := redis.UserFollowers.Exists(redis.Ctx, strconv.FormatInt(targetId, 10)).Result()

	if err2 != nil {
		log.Printf("Check Redis existence failed: %v", err2)
	}

	if keyCnt2 <= 0 {
		//获取target的粉丝，直接刷新，关注时刷新target的粉丝
		userFollowersId, _, err := followDao.GetFollowersInfo(targetId)
		if err != nil {
			log.Printf("GetFollowersInfo failed: %v", err)
			return
		}
		ImportToRDBFollower(targetId, userFollowersId)
	}

	redis.UserFollowers.SAdd(redis.Ctx, strconv.FormatInt(targetId, 10), userId)

	// 进行好友的判定，本接口实现user对target的关注，若此时target也关注了user，进行friend数据库的记录追加
	// user的好友有target，target的好友有user
	if flag, _ := followService.CheckIsFollowing(targetId, userId); flag {
		// 尝试给friend数据库追加user的好友有target的记录
		keyCnt3, err3 := redis.UserFriends.Exists(redis.Ctx, strconv.FormatInt(userId, 10)).Result()

		if err3 != nil {
			log.Printf("Check Redis existence failed: %v", err3)
		}
		if keyCnt3 <= 0 {
			userFriendsId1, _, err := followDao.GetFriendsInfo(userId)
			if err != nil {
				log.Printf("GetFriendsInfo failed: %v", err)
				return
			}
			ImportToRDBFriend(userId, userFriendsId1)
		}

		redis.UserFriends.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), targetId)

		// 尝试给friend数据库追加target的好友有user的记录
		keyCnt4, err4 := redis.UserFriends.Exists(redis.Ctx, strconv.FormatInt(targetId, 10)).Result()

		if err4 != nil {
			log.Printf("Check Redis existence failed: %v", err4)
		}
		if keyCnt4 <= 0 {
			//获取target的好友，直接刷新，关注时刷新target的好友
			userFriendsId2, _, err := followDao.GetFriendsInfo(targetId)
			if err != nil {
				log.Printf("GetFriendsInfo failed: %v", err)
				return
			}
			ImportToRDBFriend(targetId, userFriendsId2)
		}

		redis.UserFriends.SAdd(redis.Ctx, strconv.FormatInt(targetId, 10), userId)
	}
}

/*
	取关业务
*/

// CancelFollowAction 取关操作的业务
func (followService *FollowServiceImp) CancelFollowAction(userId int64, targetId int64) (bool, error) {

	followDao := repository.NewFollowDaoInstance()
	follow, err := followDao.FindEverFollowing(userId, targetId)
	// 寻找 SQL 出错。
	if nil != err {
		return false, err
	}
	// 曾经关注过，只需要update一下cancel即可。
	if nil != follow {
		_, err := followDao.UpdateFollowRelation(userId, targetId, int8(0))
		if err != nil {
			log.Printf("UpdateFollowRelation (cancel) failed: %v", err)
			return false, err
		}

		// 获取取关的消息队列
		followDelMQ := rabbitmq.SimpleFollowDelMQ
		err = followDelMQ.PublishSimpleFollow(fmt.Sprintf("%d-%d-%s", userId, targetId, "update"))
		if err != nil {
			log.Printf("Publish cancel follow message failed: %v", err)
		}

		DelToRDBWhenCancelFollow(userId, targetId)
		return true, nil
	}
	// 没有关注关系
	return false, nil
}
func DelToRDBWhenCancelFollow(userId int64, targetId int64) {
	// 当a取关b时，redis的三个关注数据库会有以下操作
	//代码使用了 redis.UserFollowings.SRem 方法，该方法用于从 Redis 的集合中移除一个或多个成员。
	redis.UserFollowings.SRem(redis.Ctx, strconv.FormatInt(userId, 10), targetId)

	redis.UserFollowers.SRem(redis.Ctx, strconv.FormatInt(targetId, 10), userId)

	// a取关b，如果a和b属于互关的用户，则两者的互关记录都会删除
	redis.UserFriends.SRem(redis.Ctx, strconv.FormatInt(userId, 10), targetId)
	redis.UserFriends.SRem(redis.Ctx, strconv.FormatInt(targetId, 10), userId)
}

/*
	获取关注列表业务
*/

// GetFollowingsByRedis 从redis获取登陆用户关注列表
func GetFollowingsByRedis(userId int64) ([]int64, int64, error) {
	followDao := repository.NewFollowDaoInstance()
	// 判定键是否存在
	keyCnt, err := redis.UserFollowings.Exists(redis.Ctx, strconv.FormatInt(userId, 10)).Result()

	if err != nil {
		log.Printf("Check Redis existence failed: %v", err)
		return nil, 0, err
	}

	// 若键存在，获取缓存数据后返回
	if keyCnt > 0 {
		ids := redis.UserFollowings.SMembers(redis.Ctx, strconv.FormatInt(userId, 10)).Val()
		idsInt64, err := convertToInt64Array(ids)
		if err != nil {
			log.Printf("ConvertToInt64Array failed: %v", err)
			return nil, 0, err
		}
		return idsInt64, int64(len(idsInt64)), nil
	} else {
		// 键不存在，获取数据库数据，更新缓存并返回
		userFollowingsId, userFollowingsCnt, err1 := followDao.GetFollowingsInfo(userId)
		if err1 != nil {
			log.Printf("GetFollowingsInfo failed: %v", err1)
			return nil, 0, err1
		}
		ImportToRDBFollowing(userId, userFollowingsId)
		return userFollowingsId, userFollowingsCnt, nil
	}

}

// GetFollowings 获取正在关注的用户详情列表业务
func (followService *FollowServiceImp) GetFollowings(userId int64) ([]User, error) {
	// 调用集成redis的关注用户获取接口获取关注用户id和关注用户数量
	userFollowingsId, userFollowingsCnt, err := GetFollowingsByRedis(userId)
	if nil != err {
		log.Printf("GetFollowingsByRedis failed: %v", err)
	}

	// 根据关注用户数量创建空用户结构体数组
	userFollowings := make([]User, userFollowingsCnt)

	// 传入buildtype调用用户构建函数构建关注用户数组
	err1 := followService.BuildUser(userId, userFollowings, userFollowingsId, 0)

	if nil != err1 {
		log.Printf("BuildUser failed: %v", err1)
	}

	return userFollowings, nil
}

/*
	获取粉丝列表业务
*/

// GetFollowersByRedis 从redis中获取用户粉丝列表
func GetFollowersByRedis(userId int64) ([]int64, int64, error) {
	followDao := repository.NewFollowDaoInstance()
	keyCnt, err := redis.UserFollowers.Exists(redis.Ctx, strconv.FormatInt(userId, 10)).Result()

	if err != nil {
		log.Printf("Check Redis existence failed: %v", err)
		return nil, 0, err
	}

	if keyCnt > 0 {
		// 键存在，获取键中集合元素
		ids := redis.UserFollowers.SMembers(redis.Ctx, strconv.FormatInt(userId, 10)).Val()
		idsInt64, err := convertToInt64Array(ids)
		if err != nil {
			log.Printf("ConvertToInt64Array failed: %v", err)
			return nil, 0, err
		}
		return idsInt64, int64(len(idsInt64)), nil
	} else {
		// 键不存在，获取数据库数据更新至redis，返回数据库所获取数据
		userFollowersId, userFollowersCnt, err1 := followDao.GetFollowersInfo(userId)
		if err1 != nil {
			log.Printf("GetFollowersInfo failed: %v", err1)
			return nil, 0, err1
		}
		ImportToRDBFollower(userId, userFollowersId)
		return userFollowersId, userFollowersCnt, nil
	}

}

// GetFollowers 获取粉丝详情列表业务
func (followService *FollowServiceImp) GetFollowers(userId int64) ([]User, error) {
	// 调用集成redis的粉丝获取接口获取粉丝id和粉丝数量
	userFollowersId, userFollowersCnt, err := GetFollowersByRedis(userId)

	if nil != err {
		log.Printf("GetFollowersByRedis failed: %v", err)
	}

	// 根据粉丝数量创建空用户结构体数组
	userFollowers := make([]User, userFollowersCnt)

	// 传入buildtype调用用户构建函数构建粉丝数组
	err1 := followService.BuildUser(userId, userFollowers, userFollowersId, 1)

	if nil != err1 {
		log.Printf("BuildUser failed: %v", err1)
	}

	return userFollowers, nil

}

/*
	获取用户好友列表业务
*/

// 从redis中获取好友信息
func GetFriendsByRedis(userId int64) ([]int64, int64, error) {
	followDao := repository.NewFollowDaoInstance()
	keyCnt, err := redis.UserFriends.Exists(redis.Ctx, strconv.FormatInt(userId, 10)).Result()

	if err != nil {
		log.Printf("Check Redis existence failed: %v", err)
		return nil, 0, err
	}

	if keyCnt > 0 {
		// 键存在，获取键中集合元素
		ids := redis.UserFriends.SMembers(redis.Ctx, strconv.FormatInt(userId, 10)).Val()
		idsInt64, err := convertToInt64Array(ids)
		if err != nil {
			log.Printf("ConvertToInt64Array failed: %v", err)
			return nil, 0, err
		}

		return idsInt64, int64(len(idsInt64)), nil

	} else {
		// 键不存在，获取数据库数据更新至redis，返回数据库所获取数据
		userFriendsId, userFriendsCnt, err1 := followDao.GetFriendsInfo(userId)
		if err1 != nil {
			log.Printf("GetFriendsInfo failed: %v", err1)
			return nil, 0, err1
		}
		ImportToRDBFriend(userId, userFriendsId)

		return userFriendsId, userFriendsCnt, nil
	}

}

// GetFriends 获取用户好友列表（附带与其最新聊天记录）
func (followService *FollowServiceImp) GetFriends(userId int64) ([]FriendUser, error) {
	// 调用集成redis的好友获取接口获取好友id和好友数量
	userFriendId, userFriendCnt, err := GetFriendsByRedis(userId)

	if nil != err {
		log.Printf("GetFriendsByRedis failed: %v", err)
	}

	// 使用好友数量创建空好友结构体数组
	userFriends := make([]FriendUser, userFriendCnt)

	// 调用好友构建函数构建好友数组
	err1 := followService.BuildFriendUser(userId, userFriends, userFriendId)

	if err1 != nil {
		log.Printf("BuildFriendUser failed: %v", err1)
	}

	return userFriends, nil
}

/*
	对外提供服务之返回登陆用户的关注用户数量
*/

// GetFollowingCnt 加入redis 根据用户id查询关注数
func (followService *FollowServiceImp) GetFollowingCnt(userId int64) (int64, error) {
	followDao := repository.NewFollowDaoInstance()

	keyCnt, err := redis.UserFollowings.Exists(redis.Ctx, strconv.FormatInt(userId, 10)).Result()

	if err != nil {
		log.Printf("Check Redis existence failed: %v", err)
		return 0, err
	}

	if keyCnt > 0 {
		// 键存在，获取键中集合元素个数
		cnt, err2 := redis.UserFollowings.SCard(redis.Ctx, strconv.FormatInt(userId, 10)).Result()
		if err2 != nil {
			log.Printf("Get Redis count failed: %v", err2)
			return 0, err2
		}
		if cnt < 100 { // 数据小的场景下 查询耗时短 数据差异明显 缓存需及时更新
			redis.UserFollowings.Expire(redis.Ctx, strconv.Itoa(int(userId)), time.Duration(5+int(time.Second)*rand.Intn(5)))
		} else if cnt < 1000 {
			redis.UserFollowings.Expire(redis.Ctx, strconv.Itoa(int(userId)), time.Duration(int(time.Second)*rand.Intn(5))+time.Minute)
		} else { // 数据大的场景 查询耗时长 数据变化不明显 可以更长时间更新一次
			redis.UserFollowings.Expire(redis.Ctx, strconv.Itoa(int(userId)), CacheTimeGenerator())
		}
		return cnt, nil

	} else {
		// 键不存在，获取数据库数据更新至redis，返回数据库所获取数据
		ids, _, err1 := followDao.GetFollowingsInfo(userId)
		if err1 != nil {
			log.Printf("GetFollowingsInfo failed: %v", err1)
			return 0, err1
		}

		ImportToRDBFollowing(userId, ids)

		return int64(len(ids)), nil
	}

}

/*
	对外提供服务之返回登陆用户的粉丝用户数量
*/

// GetFollowerCnt 根据用户id查询粉丝数
func (followService *FollowServiceImp) GetFollowerCnt(userId int64) (int64, error) {
	followDao := repository.NewFollowDaoInstance()

	keyCnt, err := redis.UserFollowers.Exists(redis.Ctx, strconv.FormatInt(userId, 10)).Result()

	if err != nil {
		log.Printf("Check Redis existence failed: %v", err)
		return 0, err
	}

	if keyCnt > 0 {
		// 键存在，获取键中集合元素个数
		cnt, err2 := redis.UserFollowers.SCard(redis.Ctx, strconv.Itoa(int(userId))).Result()

		if err2 != nil {
			log.Printf("Get Redis count failed: %v", err2)
			return 0, err2
		}
		if cnt < 100 { // 数据小的场景下 查询耗时短 数据差异明显 缓存需及时更新
			redis.UserFollowers.Expire(redis.Ctx, strconv.Itoa(int(userId)), time.Duration(5+int(time.Second)*rand.Intn(5)))
		} else if cnt < 1000 {
			redis.UserFollowers.Expire(redis.Ctx, strconv.Itoa(int(userId)), time.Duration(int(time.Second)*rand.Intn(5))+time.Minute)
		} else { // 数据大的场景 查询耗时长 数据变化不明显 可以更长时间更新一次
			redis.UserFollowers.Expire(redis.Ctx, strconv.Itoa(int(userId)), CacheTimeGenerator())
		}
		return cnt, nil

	} else {
		// 键不存在，获取数据库数据更新至redis，返回数据库所获取数据
		ids, _, err1 := followDao.GetFollowersInfo(userId)

		if err1 != nil {
			log.Printf("GetFollowersInfo failed: %v", err1)
			return 0, err1
		}

		ImportToRDBFollower(userId, ids)

		return int64(len(ids)), nil
	}

}

/*
	对外提供服务之返回登陆用户是否关注目标用户的布尔值
*/

// CheckIsFollowing 判断当前登录用户是否关注了目标用户
func (followService *FollowServiceImp) CheckIsFollowing(userId int64, targetId int64) (bool, error) {
	followDao := repository.NewFollowDaoInstance()

	keyCnt, err := redis.UserFollowings.Exists(redis.Ctx, strconv.FormatInt(userId, 10)).Result()

	if err != nil {
		log.Printf("Check Redis existence failed: %v", err)
		return false, err
	}

	if keyCnt > 0 {
		// 键存在判断是否存在userId和targetId键值对
		flag, err3 := redis.UserFollowings.SIsMember(redis.Ctx, strconv.Itoa(int(userId)), targetId).Result()

		if err3 != nil {
			log.Printf("Check Redis membership failed: %v", err3)
			return false, err3
		}

		if flag {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		// 键不存在，获取数据库数据更新至redis中，使用dao层方法判断是否有关注关系
		ids, _, err1 := followDao.GetFollowingsInfo(userId)

		if err1 != nil {
			log.Printf("GetFollowingsInfo failed: %v", err1)
			return false, err1
		}

		ImportToRDBFollowing(userId, ids)

		isFollow, err2 := followDao.FindFollowRelation(userId, targetId)

		if err2 != nil {
			log.Printf("FindFollowRelation failed: %v", err2)
			return false, err2
		}

		return isFollow, nil
	}

}

/*
	提供目标用户id和对应的id列表导入到redis中的方法，一般用在更新失效键的逻辑中
*/

// ImportToRDBFollowing 将登陆用户的关注id列表导入到following数据库中
func ImportToRDBFollowing(userId int64, ids []int64) {
	// 将传入的userId及其关注用户id更新至redis中
	for _, id := range ids {
		redis.UserFollowings.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), int(id))
	}
	//通过 redis.UserFollowings.Expire 方法设置 Redis 数据库中存储用户关注关系的集合的过期时间，过期时间由 CacheTimeGenerator() 函数生成。
	redis.UserFollowings.Expire(redis.Ctx, strconv.FormatInt(userId, 10), CacheTimeGenerator())
}

// ImportToRDBFollower 将登陆用户的关注id列表导入到follower数据库中
func ImportToRDBFollower(userId int64, ids []int64) {
	// 将传入的userId及其粉丝id更新至redis中
	for _, id := range ids {
		redis.UserFollowers.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), int(id))
	}

	redis.UserFollowers.Expire(redis.Ctx, strconv.FormatInt(userId, 10), CacheTimeGenerator())
}

func ImportToRDBFriend(userId int64, ids []int64) {
	// 将传入的userId及其好友id更新至redis中
	for _, id := range ids {
		redis.UserFriends.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), int(id))
	}

	redis.UserFriends.Expire(redis.Ctx, strconv.FormatInt(userId, 10), CacheTimeGenerator())
}

/*
	将返回关注用户、返回粉丝用户、返回好友用户中的构建用户的逻辑独立出来
	注： builduser方法根据传入的buildtype决定是构建关注用户还是粉丝用户
*/

// BatchGetUserNames 批量获取用户名，减少数据库查询次数
func (followService *FollowServiceImp) BatchGetUserNames(ids []int64) (map[int64]string, error) {
	if len(ids) == 0 {
		return make(map[int64]string), nil
	}

	u := dao.User
	userList, err := u.Where(u.ID.In(ids...)).Find()
	if err != nil {
		log.Printf("BatchGetUserNames failed: %v", err)
		return nil, err
	}

	nameMap := make(map[int64]string, len(userList))
	for _, user := range userList {
		nameMap[user.ID] = user.Name
	}

	return nameMap, nil
}

// BatchGetFollowingCounts 批量获取关注数
func (followService *FollowServiceImp) BatchGetFollowingCounts(ids []int64) (map[int64]int64, error) {
	if len(ids) == 0 {
		return make(map[int64]int64), nil
	}

	followDao := repository.NewFollowDaoInstance()
	countMap := make(map[int64]int64, len(ids))

	for _, id := range ids {
		keyCnt, err := redis.UserFollowings.Exists(redis.Ctx, strconv.FormatInt(id, 10)).Result()
		if err != nil {
			log.Printf("Check Redis existence failed: %v", err)
			continue
		}

		if keyCnt > 0 {
			cnt, _ := redis.UserFollowings.SCard(redis.Ctx, strconv.FormatInt(id, 10)).Result()
			countMap[id] = cnt
		} else {
			cnt, _ := followDao.GetFollowingCnt(id)
			countMap[id] = cnt
		}
	}

	return countMap, nil
}

// BatchGetFollowerCounts 批量获取粉丝数
func (followService *FollowServiceImp) BatchGetFollowerCounts(ids []int64) (map[int64]int64, error) {
	if len(ids) == 0 {
		return make(map[int64]int64), nil
	}

	followDao := repository.NewFollowDaoInstance()
	countMap := make(map[int64]int64, len(ids))

	for _, id := range ids {
		keyCnt, err := redis.UserFollowers.Exists(redis.Ctx, strconv.FormatInt(id, 10)).Result()
		if err != nil {
			log.Printf("Check Redis existence failed: %v", err)
			continue
		}

		if keyCnt > 0 {
			cnt, _ := redis.UserFollowers.SCard(redis.Ctx, strconv.FormatInt(id, 10)).Result()
			countMap[id] = cnt
		} else {
			cnt, _ := followDao.GetFollowerCnt(id)
			countMap[id] = cnt
		}
	}

	return countMap, nil
}

// BatchCheckIsFollowing 批量检查用户是否关注了目标用户
func (followService *FollowServiceImp) BatchCheckIsFollowing(userId int64, targetIds []int64) (map[int64]bool, error) {
	if len(targetIds) == 0 {
		return make(map[int64]bool), nil
	}

	resultMap := make(map[int64]bool, len(targetIds))

	keyCnt, err := redis.UserFollowings.Exists(redis.Ctx, strconv.FormatInt(userId, 10)).Result()
	if err != nil {
		log.Printf("Check Redis existence failed: %v", err)
	}

	if keyCnt > 0 {
		for _, targetId := range targetIds {
			flag, _ := redis.UserFollowings.SIsMember(redis.Ctx, strconv.FormatInt(userId, 10), targetId).Result()
			resultMap[targetId] = flag
		}
	} else {
		followDao := repository.NewFollowDaoInstance()
		for _, targetId := range targetIds {
			flag, _ := followDao.FindFollowRelation(userId, targetId)
			resultMap[targetId] = flag
		}
	}

	return resultMap, nil
}

// BuildUser 根据传入的id列表和空user数组，构建业务所需user数组并返回
func (followService *FollowServiceImp) BuildUser(userId int64, users []User, ids []int64, buildtype int) error {
	if len(ids) == 0 {
		return nil
	}

	// 批量获取所有用户信息
	userNameMap, err := followService.BatchGetUserNames(ids)
	if err != nil {
		log.Printf("BatchGetUserNames failed: %v", err)
		return err
	}

	followingCountMap, err := followService.BatchGetFollowingCounts(ids)
	if err != nil {
		log.Printf("BatchGetFollowingCounts failed: %v", err)
		return err
	}

	followerCountMap, err := followService.BatchGetFollowerCounts(ids)
	if err != nil {
		log.Printf("BatchGetFollowerCounts failed: %v", err)
		return err
	}

	var isFollowMap map[int64]bool
	if buildtype == 1 {
		isFollowMap, err = followService.BatchCheckIsFollowing(userId, ids)
		if err != nil {
			log.Printf("BatchCheckIsFollowing failed: %v", err)
			return err
		}
	}

	for i, id := range ids {
		users[i].Id = id
		users[i].Name = userNameMap[id]

		if users[i].Name == "" {
			users[i].Name = "未知用户"
		}

		users[i].FollowCount = followingCountMap[id]
		users[i].FollowerCount = followerCountMap[id]

		if buildtype == 1 {
			users[i].IsFollow = isFollowMap[id]
		} else {
			users[i].IsFollow = true
		}
	}

	return nil
}

// BatchGetLatestMessages 批量量取与多个好友的最新消息
func (followService *FollowServiceImp) BatchGetLatestMessages(userId int64, friendIds []int64) (map[int64]*model.Message, error) {
	if len(friendIds) == 0 {
		return make(map[int64]*model.Message), nil
	}

	messageMap := make(map[int64]*model.Message, len(friendIds))

	for _, friendId := range friendIds {
		messageInfo, err := followService.GetLatestMessage(userId, friendId)
		if err != nil {
			log.Printf("GetLatestMessage failed for userId %d, friendId %d: %v", userId, friendId, err)
			continue
		}
		messageMap[friendId] = messageInfo
	}

	return messageMap, nil
}

// BuildFriendUser 根据传入的id列表和空frienduser数组，构建业务所需frienduser数组并返回
func (followService *FollowServiceImp) BuildFriendUser(userId int64, friendUsers []FriendUser, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	// 批量获取所有用户信息
	userNameMap, err := followService.BatchGetUserNames(ids)
	if err != nil {
		log.Printf("BatchGetUserNames failed: %v", err)
		return err
	}

	followingCountMap, err := followService.BatchGetFollowingCounts(ids)
	if err != nil {
		log.Printf("BatchGetFollowingCounts failed: %v", err)
		return err
	}

	followerCountMap, err := followService.BatchGetFollowerCounts(ids)
	if err != nil {
		log.Printf("BatchGetFollowerCounts failed: %v", err)
		return err
	}

	messageMap, err := followService.BatchGetLatestMessages(userId, ids)
	if err != nil {
		log.Printf("BatchGetLatestMessages failed: %v", err)
	}

	defaultAvatar := viper.GetString("settings.oss.avatar")

	for i, id := range ids {
		friendUsers[i].Id = id
		friendUsers[i].Name = userNameMap[id]

		if friendUsers[i].Name == "" {
			friendUsers[i].Name = "未知用户"
		}

		friendUsers[i].FollowCount = followingCountMap[id]
		friendUsers[i].FollowerCount = followerCountMap[id]
		friendUsers[i].IsFollow = true
		friendUsers[i].Avatar = defaultAvatar

		if messageInfo, exists := messageMap[id]; exists && messageInfo != nil {
			friendUsers[i].MsgContent = messageInfo.Content
			friendUsers[i].MsgType = messageInfo.ActionType
		}
	}

	return nil
}
