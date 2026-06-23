package utils

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var ctx = context.Background()

func AcquireLock(client *redis.Client, lockKey string, timeout time.Duration) bool {
	isLocked, err := client.SetNX(ctx, lockKey, "locked", timeout).Result()
	if err != nil {
		fmt.Println(err)
	}
	return isLocked
}

func ReleaseLock(client *redis.Client, lockKey string) {
	_, err := client.Del(ctx, lockKey).Result()
	if err != nil {
		fmt.Println(err)
	}
}

//func likePost(client *redis.Client, postID string) (int64, error) {
//	lockKey := "post:" + postID + ":lock"
//	timeout := 5 * time.Second
//
//	if !acquireLock(client, lockKey, timeout) {
//		return 0, fmt.Errorf("could not acquire lock for post %s", postID)
//	}
//	defer releaseLock(client, lockKey)
//
//	// 假设redis中存储了post的点赞数量
//	likeCountKey := "post:" + postID + ":likes"
//	newLikeCount, err := client.Incr(ctx, likeCountKey).Result()
//	if err != nil {
//		return 0, err
//	}
//	return newLikeCount, nil
//}
//
//func main() {
//	client := redis.NewClient(&redis.Options{
//		Addr:     "localhost:6379",
//		Password: "", // 无密码时留空
//		DB:       0,  // 默认数据库
//	})
//
//	postID := "12345"
//	likeCount, err := likePost(client, postID)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("Post %s has been liked %d times.\n", postID, likeCount)
//}
