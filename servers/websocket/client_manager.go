/**
 * Created by GoLand.
 * User: nankai
 * Date: 2019-07-25
 * Time: 16:24
 */

package websocket

import (
	"fmt"
	"gowebsocket/helper"
	"gowebsocket/lib/cache"
	"gowebsocket/models"
	"sync"
	"time"
)

// connection management
type ClientManager struct {
	Clients     map[*Client]bool   // all connections
	ClientsLock sync.RWMutex       // Read-write lock
	Users       map[string]*Client // logged in users // appId+uuid
	UserLock    sync.RWMutex       // Read-write lock
	Register    chan *Client       // Connection connection processing
	Login       chan *login        // User login processing
	Unregister  chan *Client       // disconnect handler
	Broadcast   chan []byte        // Broadcast to send data to all members
}

func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		Clients:    make(map[*Client]bool),
		Users:      make(map[string]*Client),
		Register:   make(chan *Client, 1000),
		Login:      make(chan *login, 1000),
		Unregister: make(chan *Client, 1000),
		Broadcast:  make(chan []byte, 1000),
	}

	return
}

// get user key
func GetUserKey(appId uint32, userId string) (key string) {
	key = fmt.Sprintf("%d_%s", appId, userId)

	return
}

/**************************** manager ********************** ****************/

func (manager *ClientManager) InClient(client *Client) (ok bool) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	// The connection exists, after adding
	_, ok = manager.Clients[client]

	return
}

// GetClients
func (manager *ClientManager) GetClients() (clients map[*Client]bool) {

	clients = make(map[*Client]bool)

	manager.ClientsRange(func(client *Client, value bool) (result bool) {
		clients[client] = value

		return true
	})

	return
}

// traverse
func (manager *ClientManager) ClientsRange(f func(client *Client, value bool) (result bool)) {

	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	for key, value := range manager.Clients {
		result := f(key, value)
		if result == false {
			return
		}
	}

	return
}

// GetClientsLen
func (manager *ClientManager) GetClientsLen() (clientsLen int) {

	clientsLen = len(manager.Clients)

	return
}

// add client
func (manager *ClientManager) AddClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	manager.Clients[client] = true
}

// delete client
func (manager *ClientManager) DelClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	if _, ok := manager.Clients[client]; ok {
		delete(manager.Clients, client)
	}
}

// Get the user's connection
func (manager *ClientManager) GetUserClient(appId uint32, userId string) (client *Client) {

	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()

	userKey := GetUserKey(appId, userId)
	if value, ok := manager.Users[userKey]; ok {
		client = value
	}

	return
}

// GetClientsLen
func (manager *ClientManager) GetUsersLen() (userLen int) {
	userLen = len(manager.Users)

	return
}

// Add user
func (manager *ClientManager) AddUsers(key string, client *Client) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()

	manager.Users[key] = client
}

// delete users
func (manager *ClientManager) DelUsers(client *Client) (result bool) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()

	key := GetUserKey(client.AppId, client.UserId)
	if value, ok := manager.Users[key]; ok {
		// Determine if it is the same user
		if value.Addr != client.Addr {

			return
		}
		delete(manager.Users, key)
		result = true
	}

	return
}

// Get the user's key
func (manager *ClientManager) GetUserKeys() (userKeys []string) {

	userKeys = make([]string, 0)
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()
	for key := range manager.Users {
		userKeys = append(userKeys, key)
	}

	return
}

// Get the user's key
func (manager *ClientManager) GetUserList(appId uint32) (userList []string) {

	userList = make([]string, 0)

	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()

	for _, v := range manager.Users {
		if v.AppId == appId {
			userList = append(userList, v.UserId)
		}
	}

	fmt.Println("GetUserList len:", len(manager.Users))

	return
}

// Get the user's key
func (manager *ClientManager) GetUserClients() (clients []*Client) {

	clients = make([]*Client, 0)
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()
	for _, v := range manager.Users {
		clients = append(clients, v)
	}

	return
}

// send data to all members (except myself)
func (manager *ClientManager) sendAll(message []byte, ignoreClient *Client) {

	clients := manager.GetUserClients()
	for _, conn := range clients {
		if conn != ignoreClient {
			conn.SendMsg(message)
		}
	}
}

// send data to all members (except myself)
func (manager *ClientManager) sendAppIdAll(message []byte, appId uint32, ignoreClient *Client) {

	clients := manager.GetUserClients()
	for _, conn := range clients {
		if conn != ignoreClient && conn.AppId == appId {
			conn.SendMsg(message)
		}
	}
}

// User establish connection event
func (manager *ClientManager) EventRegister(client *Client) {
	manager.AddClients(client)

	fmt.Println("EventRegister user established connection", client.Addr)

	// client.Send <- []byte("Connection succeeded")
}

// User login
func (manager *ClientManager) EventLogin(login *login) {

	client := login.Client
	// The connection exists, after adding
	if manager.InClient(client) {
		userKey := login.GetKey()
		manager.AddUsers(userKey, login.Client)
	}

	fmt.Println("EventLogin user login", client.Addr, login.AppId, login.UserId)

	orderId := helper.GetOrderIdTime()
	SendUserMessageAll(login.AppId, login.UserId, orderId, models.MessageCmdEnter, "Hello~")
}

// user disconnects
func (manager *ClientManager) EventUnregister(client *Client) {
	manager.DelClients(client)

	// delete user connection
	deleteResult := manager.DelUsers(client)
	if deleteResult == false {
		// not currently connected client

		return
	}

	// clear redis login data
	userOnline, err := cache.GetUserOnlineInfo(client.GetKey())
	if err == nil {
		userOnline.LogOut()
		cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	}

	// close chan
	// close(client.Send)

	fmt.Println("EventUnregister user disconnected", client.Addr, client.AppId, client.UserId)

	if client.UserId != "" {
		orderId := helper.GetOrderIdTime()
		SendUserMessageAll(client.AppId, client.UserId, orderId, models.MessageCmdExit, "User has left~")
	}
}

// pipe handler
func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.Register:
			// create connection event
			manager.EventRegister(conn)

		case login := <-manager.Login:
			// User login
			manager.EventLogin(login)

		case conn := <-manager.Unregister:
			// disconnect event
			manager.EventUnregister(conn)

		case message := <-manager.Broadcast:
			// broadcast event
			clients := manager.GetClients()
			for conn := range clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
				}
			}
		}
	}
}

/**************************** manager info ******************** ********************/
// get manager information
func GetManagerInfo(isDebug string) (managerInfo map[string]interface{}) {
	managerInfo = make(map[string]interface{})

	managerInfo["clientsLen"] = clientManager.GetClientsLen()        // Number of client connections
	managerInfo["usersLen"] = clientManager.GetUsersLen()            // Number of logged in users
	managerInfo["chanRegisterLen"] = len(clientManager.Register)     // number of unhandled connection events
	managerInfo["chanLoginLen"] = len(clientManager.Login)           // Number of unhandled login events
	managerInfo["chanUnregisterLen"] = len(clientManager.Unregister) // Number of unhandled logout events
	managerInfo["chanBroadcastLen"] = len(clientManager.Broadcast)   // number of unhandled broadcast events

	if isDebug == "true" {
		addrList := make([]string, 0)
		clientManager.ClientsRange(func(client *Client, value bool) (result bool) {
			addrList = append(addrList, client.Addr)

			return true
		})

		users := clientManager.GetUserKeys()

		managerInfo["clients"] = addrList // client list
		managerInfo["users"] = users      // list of logged in users
	}

	return
}

// Get the connection the user is on
func GetUserClient(appId uint32, userId string) (client *Client) {
	client = clientManager.GetUserClient(appId, userId)

	return
}

// Periodically clean up the timeout connection
func ClearTimeoutConnections() {
	currentTime := uint64(time.Now().Unix())

	clients := clientManager.GetClients()
	for client := range clients {
		if client.IsHeartbeatTimeout(currentTime) {
			fmt.Println("Heartbeat time timeout, close the connection", client.Addr, client.UserId, client.LoginTime, client.HeartbeatTime)

			client.Socket.Close()
		}
	}
}

// get all users
func GetUserList(appId uint32) (userList []string) {
	fmt.Println("Get all users", appId)

	userList = clientManager.GetUserList(appId)

	return
}

// all broadcast
func AllSendMessages(appId uint32, userId string, data string) {
	fmt.Println("All staff broadcast", appId, userId, data)

	ignoreClient := clientManager.GetUserClient(appId, userId)
	clientManager.sendAppIdAll([]byte(data), appId, ignoreClient)
}
