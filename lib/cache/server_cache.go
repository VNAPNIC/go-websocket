/**
* Created by GoLand.
* User: nankai
* Date: 2019-08-03
* Time: 15:23
 */

package cache

import (
	"encoding/json"
	"fmt"
	"gowebsocket/lib/redislib"
	"gowebsocket/models"
	"strconv"
)

const (
	serversHashKey       = "acc:hash:servers" // all servers
	serversHashCacheTime = 2 * 60 * 60        // key expiration time
	serversHashTimeout   = 3 * 60             // timeout
)

func getServersHashKey() (key string) {
	key = fmt.Sprintf("%s", serversHashKey)

	return
}

// Set server information
func SetServerInfo(server *models.Server, currentTime uint64) (err error) {
	key := getServersHashKey()

	value := fmt.Sprintf("%d", currentTime)

	redisClient := redislib.GetClient()
	number, err := redisClient.Do("hSet", key, server.String(), value).Int()
	if err != nil {
		fmt.Println("SetServerInfo", key, number, err)

		return
	}

	if number != 1 {

		return
	}

	redisClient.Do("Expire", key, serversHashCacheTime)

	return
}

// Offline server information
func DelServerInfo(server *models.Server) (err error) {
	key := getServersHashKey()
	redisClient := redislib.GetClient()
	number, err := redisClient.Do("hDel", key, server.String()).Int()
	if err != nil {
		fmt.Println("DelServerInfo", key, number, err)

		return
	}

	if number != 1 {

		return
	}

	redisClient.Do("Expire", key, serversHashCacheTime)

	return
}

func GetServerAll(currentTime uint64) (servers []*models.Server, err error) {

	servers = make([]*models.Server, 0)
	key := getServersHashKey()

	redisClient := redislib.GetClient()

	val, err := redisClient.Do("hGetAll", key).Result()

	valByte, _ := json.Marshal(val)
	fmt.Println("GetServerAll", key, string(valByte))

	serverMap, err := redisClient.HGetAll(key).Result()
	if err != nil {
		fmt.Println("SetServerInfo", key, err)

		return
	}

	for key, value := range serverMap {
		valueUint64, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			fmt.Println("GetServerAll", key, err)

			return nil, err
		}

		// time out
		if valueUint64+serversHashTimeout <= currentTime {
			continue
		}

		server, err := models.StringToServer(key)
		if err != nil {
			fmt.Println("GetServerAll", key, err)

			return nil, err
		}

		servers = append(servers, server)
	}

	return
}
