/**
* Created by GoLand.
* User: nankai
* Date: 2019-07-25
* Time: 12:11
 */

package common

const (
	OK                 = 200  // Success
	NotLoggedIn        = 1000 // not logged in
	ParameterIllegal   = 1001 // parameter is invalid
	UnauthorizedUserId = 1002 // Illegal UserId
	Unauthorized       = 1003 // Unauthorized
	ServerError        = 1004 // system error
	NotData            = 1005 // no data
	ModelAddError      = 1006 // add error
	ModelDeleteError   = 1007 // delete error
	ModelStoreError    = 1008 // store error
	OperationFailure   = 1009 // Operation failed
	RoutingNotExist    = 1010 // route does not exist
)

// Get error information based on error code
func GetErrorMessage(code uint32, message string) string {
	var codeMessage string
	codeMap := map[uint32]string{
		OK:                 "Success",
		NotLoggedIn:        "Not logged in",
		ParameterIllegal:   "Parameter is invalid",
		UnauthorizedUserId: "Illegal UserId",
		Unauthorized:       "Unauthorized",
		NotData:            "No data",
		ServerError:        "System error",
		ModelAddError:      "Add Error",
		ModelDeleteError:   "Delete Error",
		ModelStoreError:    "Storage error",
		OperationFailure:   "Operation failed",
		RoutingNotExist:    "Routing does not exist",
	}

	if message == "" {
		if value, ok := codeMap[code]; ok {
			// exist
			codeMessage = value
		} else {
			codeMessage = "undefined error type!"
		}
	} else {
		codeMessage = message
	}

	return codeMessage
}
