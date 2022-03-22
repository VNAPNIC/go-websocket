/**
 * Created by GoLand.
 * User: nankai
 * Date: 2019-07-25
 * Time: 17:28
 */

package cache

import (
	"encoding/json"
	"fmt"
	"go-websocket/lib/redislib"
	"go-websocket/models"

	"github.com/go-redis/redis"
)

const (
	userOnlinePrefix    = "acc:user:online:" // User online status
	userOnlineCacheTime = 24 * 60 * 60
)

/************************ Check if the user is online ************************/
func getUserOnlineKey(userKey string) (key string) {
	key = fmt.Sprintf("%s%s", userOnlinePrefix, userKey)

	return
}

func GetUserOnlineInfo(userKey string) (userOnline *models.UserOnline, err error) {
	redisClient := redislib.GetClient()

	key := getUserOnlineKey(userKey)

	data, err := redisClient.Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("GetUserOnlineInfo", userKey, err)

			return
		}

		fmt.Println("GetUserOnlineInfo", userKey, err)

		return
	}

	userOnline = &models.UserOnline{}
	err = json.Unmarshal(data, userOnline)
	if err != nil {
		fmt.Println("Get user online data json Unmarshal", userKey, err)

		return
	}

	fmt.Println("Get user online data", userKey, "time", userOnline.LoginTime, userOnline.HeartbeatTime, "AccIp", userOnline.AccIp, userOnline.IsLogoff)

	return
}

// set user online data
func SetUserOnlineInfo(userKey string, userOnline *models.UserOnline) (err error) {

	redisClient := redislib.GetClient()
	key := getUserOnlineKey(userKey)

	valueByte, err := json.Marshal(userOnline)
	if err != nil {
		fmt.Println("Set user online data json Marshal", key, err)

		return
	}

	_, err = redisClient.Do("setEx", key, userOnlineCacheTime, string(valueByte)).Result()
	if err != nil {
		fmt.Println("Set user online data", key, err)

		return
	}

	return
}
