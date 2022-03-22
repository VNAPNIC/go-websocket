/**
 * Created by GoLand.
 * User: nankai
 * Date: 2019-07-26
 * Time: 09:18
 */

package cache

import (
	"fmt"
	"gowebsocket/lib/redislib"
)

const (
	submitAgainPrefix = "acc:submit:again:" // data will not be submitted repeatedly
)

/************************ Query whether the data has been processed ************************ */

// Get data and submit to remove key
func getSubmitAgainKey(from string, value string) (key string) {
	key = fmt.Sprintf("%s%s:%s", submitAgainPrefix, from, value)

	return
}

// repeated submit
// return true: repeat submission false: first submission
func submitAgain(from string, second int, value string) (isSubmitAgain bool) {

	// default repeat submission
	isSubmitAgain = true
	key := getSubmitAgainKey(from, value)

	redisClient := redislib.GetClient()
	number, err := redisClient.Do("setNx", key, "1").Int()
	if err != nil {
		fmt.Println("submitAgain", key, number, err)

		return
	}

	if number != 1 {

		return
	}
	// first commit
	isSubmitAgain = false

	redisClient.Do("Expire", key, second)

	return

}

// Seq repeat submission
func SeqDuplicates(seq string) (result bool) {
	result = submitAgain("seq", 12*60*60, seq)

	return
}
