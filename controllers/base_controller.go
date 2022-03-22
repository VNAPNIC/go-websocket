/**
* Created by GoLand.
* User: nankai
* Date: 2019-07-25
* Time: 12:11
 */

package controllers

import (
	"go-websocket/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseController struct {
	gin.Context
}

// Get all requests resolved to map
func Response(c *gin.Context, code uint32, msg string, data map[string]interface{}) {
	message := common.Response(code, msg, data)

	// allow cross domain
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Origin", "*")                                       // This is to allow access to all domains
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") // All cross-domain request methods supported by the server, in order to avoid multiple 'pre-requests' check' request
	c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT , X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last -Modified,Pragma,FooBar") // Cross-domain key settings allow browsers to parse
	c.Header("Access-Control-Allow-Credentials", "true")                                                                                                                                                    // Whether cross-domain requests require cookie information The default setting is true
	c.Set("content-type", "application/json")                                                                                                                                                               // Set the return format to json

	c.JSON(http.StatusOK, message)

	return
}
