/**
* Created by GoLand.
* User: nankai
* Date: 2019-08-03
* Time: 16:43
 */

package grpcclient

import (
	"context"
	"errors"
	"fmt"
	"gowebsocket/common"
	"gowebsocket/models"
	"gowebsocket/protobuf"
	"time"

	"google.golang.org/grpc"
)

// rpc client
// Send a message to all users
// link::https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_client/main.go
func SendMsgAll(server *models.Server, seq string, appId uint32, userId string, cmd string, message string) (sendMsgId string, err error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(server.String(), grpc.WithInsecure())
	if err != nil {
		fmt.Println("Connection failed", server.String())

		return
	}
	defer conn.Close()

	c := protobuf.NewAccServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := protobuf.SendMsgAllReq{
		Seq:    seq,
		AppId:  appId,
		UserId: userId,
		Cms:    cmd,
		Msg:    message,
	}
	rsp, err := c.SendMsgAll(ctx, &req)
	if err != nil {
		fmt.Println("Send a message to all users", err)

		return
	}

	if rsp.GetRetCode() != common.OK {
		fmt.Println("Send a message to all users", rsp.String())
		err = errors.New(fmt.Sprintf("Failed to send message code:%d", rsp.GetRetCode()))

		return
	}

	sendMsgId = rsp.GetSendMsgId()
	fmt.Println("Send message to all users successfully:", sendMsgId)

	return
}

// get user list
// link::https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_client/main.go
func GetUserList(server *models.Server, appId uint32) (userIds []string, err error) {
	userIds = make([]string, 0)

	conn, err := grpc.Dial(server.String(), grpc.WithInsecure())
	if err != nil {
		fmt.Println("Connection failed", server.String())

		return
	}
	defer conn.Close()

	c := protobuf.NewAccServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := protobuf.GetUserListReq{
		AppId: appId,
	}
	rsp, err := c.GetUserList(ctx, &req)
	if err != nil {
		fmt.Println("get user list send request error:", err)

		return
	}

	if rsp.GetRetCode() != common.OK {
		fmt.Println("get user list return code error:", rsp.String())
		err = errors.New(fmt.Sprintf("Failed to send message code:%d", rsp.GetRetCode()))

		return
	}

	userIds = rsp.GetUserId()
	fmt.Println("Get user list success:", userIds)

	return
}

// rpc client
// send messages
// link::https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_client/main.go
func SendMsg(server *models.Server, seq string, appId uint32, userId string, cmd string, msgType string, message string) (sendMsgId string, err error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(server.String(), grpc.WithInsecure())
	if err != nil {
		fmt.Println("Connection failed", server.String())

		return
	}
	defer conn.Close()

	c := protobuf.NewAccServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := protobuf.SendMsgReq{
		Seq:     seq,
		AppId:   appId,
		UserId:  userId,
		Cms:     cmd,
		Type:    msgType,
		Msg:     message,
		IsLocal: false,
	}
	rsp, err := c.SendMsg(ctx, &req)
	if err != nil {
		fmt.Println("send messages", err)

		return
	}

	if rsp.GetRetCode() != common.OK {
		fmt.Println("send messages", rsp.String())
		err = errors.New(fmt.Sprintf("Failed to send message code:%d", rsp.GetRetCode()))

		return
	}

	sendMsgId = rsp.GetSendMsgId()
	fmt.Println("message sent successfully:", sendMsgId)

	return
}
