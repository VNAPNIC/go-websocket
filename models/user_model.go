/**
 * Created by GoLand.
 * User: nankai
 * Date: 2019-07-25
 * Time: 17:36
 */

package models

import (
	"fmt"
	"time"
)

const (
	heartbeatTimeout = 3 * 60 // User heartbeat timeout
)

// user online status
type UserOnline struct {
	AccIp         string `json:"accIp"`         // acc Ip
	AccPort       string `json:"accPort"`       // acc port
	AppId         uint32 `json:"appId"`         // appId
	UserId        string `json:"userId"`        // UserId
	ClientIp      string `json:"clientIp"`      // Client Ip
	ClientPort    string `json:"clientPort"`    // client port
	LoginTime     uint64 `json:"loginTime"`     // user's last login time
	HeartbeatTime uint64 `json:"heartbeatTime"` // user's last heartbeat time
	LogOutTime    uint64 `json:"logOutTime"`    // the time the user logged out
	Qua           string `json:"qua"`           // qua
	DeviceInfo    string `json:"deviceInfo"`    // device information
	IsLogoff      bool   `json:"isLogoff"`      // Whether to log off
}

/**********************  data processing************************* ********/

// User login
func UserLogin(accIp, accPort string, appId uint32, userId string, addr string, loginTime uint64) (userOnline *UserOnline) {

	userOnline = &UserOnline{
		AccIp:         accIp,
		AccPort:       accPort,
		AppId:         appId,
		UserId:        userId,
		ClientIp:      addr,
		LoginTime:     loginTime,
		HeartbeatTime: loginTime,
		IsLogoff:      false,
	}

	return
}

// user heartbeat
func (u *UserOnline) Heartbeat(currentTime uint64) {

	u.HeartbeatTime = currentTime
	u.IsLogoff = false

	return
}

// user logs out
func (u *UserOnline) LogOut() {

	currentTime := uint64(time.Now().Unix())
	u.LogOutTime = currentTime
	u.IsLogoff = true

	return
}

/************************ Data manipulation************************** ********/

// whether the user is online
func (u *UserOnline) IsOnline() (online bool) {
	if u.IsLogoff {

		return
	}

	currentTime := uint64(time.Now().Unix())

	if u.HeartbeatTime < (currentTime - heartbeatTimeout) {
		fmt.Println("Whether the user is online heartbeat timed out", u.AppId, u.UserId, u.HeartbeatTime)

		return
	}

	if u.IsLogoff {
		fmt.Println("Whether the user is online, the user has been offline", u.AppId, u.UserId)

		return
	}

	return true
}

// Whether the user is on this machine
func (u *UserOnline) UserIsLocal(localIp, localPort string) (result bool) {

	if u.AccIp == localIp && u.AccPort == localPort {
		result = true

		return
	}

	return
}
