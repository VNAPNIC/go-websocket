/**
* Created by GoLand.
* User: nankai
* Date: 2019-08-03
* Time: 16:43
 */

package grpcserver

import (
	"context"
	"fmt"
	"gowebsocket/common"
	"gowebsocket/models"
	"gowebsocket/protobuf"
	"gowebsocket/servers/websocket"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type server struct {
}

func setErr(rsp proto.Message, code uint32, message string) {

	message = common.GetErrorMessage(code, message)
	switch v := rsp.(type) {
	case *protobuf.QueryUsersOnlineRsp:
		v.RetCode = code
		v.ErrMsg = message
	case *protobuf.SendMsgRsp:
		v.RetCode = code
		v.ErrMsg = message
	case *protobuf.SendMsgAllRsp:
		v.RetCode = code
		v.ErrMsg = message
	case *protobuf.GetUserListRsp:
		v.RetCode = code
		v.ErrMsg = message
	default:

	}

}

// Check if user is online
func (s *server) QueryUsersOnline(c context.Context, req *protobuf.QueryUsersOnlineReq) (rsp *protobuf.QueryUsersOnlineRsp, err error) {

	fmt.Println("grpc_request Check if user is online", req.String())

	rsp = &protobuf.QueryUsersOnlineRsp{}

	online := websocket.CheckUserOnline(req.GetAppId(), req.GetUserId())

	setErr(req, common.OK, "")
	rsp.Online = online

	return rsp, nil
}

// Send a message to the local user
func (s *server) SendMsg(c context.Context, req *protobuf.SendMsgReq) (rsp *protobuf.SendMsgRsp, err error) {

	fmt.Println("grpc_request Send a message to the local user", req.String())

	rsp = &protobuf.SendMsgRsp{}

	if req.GetIsLocal() {

		// not support
		setErr(rsp, common.ParameterIllegal, "")

		return
	}

	data := models.GetMsgData(req.GetUserId(), req.GetSeq(), req.GetCms(), req.GetMsg())
	sendResults, err := websocket.SendUserMessageLocal(req.GetAppId(), req.GetUserId(), data)
	if err != nil {
		fmt.Println("system error", err)
		setErr(rsp, common.ServerError, "")

		return rsp, nil
	}

	if !sendResults {
		fmt.Println("Failed to send", err)
		setErr(rsp, common.OperationFailure, "")

		return rsp, nil
	}

	setErr(rsp, common.OK, "")

	fmt.Println("grpc_response Send a message to the local user", rsp.String())
	return
}

// Send a message to all users of the machine
func (s *server) SendMsgAll(c context.Context, req *protobuf.SendMsgAllReq) (rsp *protobuf.SendMsgAllRsp, err error) {

	fmt.Println("grpc_request Send a message to all users of the machine", req.String())

	rsp = &protobuf.SendMsgAllRsp{}

	data := models.GetMsgData(req.GetUserId(), req.GetSeq(), req.GetCms(), req.GetMsg())
	websocket.AllSendMessages(req.GetAppId(), req.GetUserId(), data)

	setErr(rsp, common.OK, "")

	fmt.Println("grpc_response Send a message to all users of the machine:", rsp.String())

	return
}

// Get a list of native users
func (s *server) GetUserList(c context.Context, req *protobuf.GetUserListReq) (rsp *protobuf.GetUserListRsp, err error) {

	fmt.Println("grpc_request Get a list of native users", req.String())

	appId := req.GetAppId()
	rsp = &protobuf.GetUserListRsp{}

	// native
	userList := websocket.GetUserList(appId)

	setErr(rsp, common.OK, "")
	rsp.UserId = userList

	fmt.Println("grpc_response Get a list of native users:", rsp.String())

	return
}

// rpc server
// link::https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_server/main.go
func Init() {

	rpcPort := viper.GetString("app.rpcPort")
	fmt.Println("rpc server start up", rpcPort)

	lis, err := net.Listen("tcp", ":"+rpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	protobuf.RegisterAccServerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
