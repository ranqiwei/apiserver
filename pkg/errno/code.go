package errno

var (
	OK                  = &Errno{Code: 0, Message: "OK."}
	InternalServerError = &Errno{Code: 10001, Message: "Internal server error."}
	ErrBind             = &Errno{Code: 10002, Message: "Error occurred while binding request to the body."}
	ErrUserNotFound     = &Errno{Code: 20102, Message: "The user was not found."}

	//common errors
	ErrValidation = &Errno{Code: 20001, Message: "Validation failed."}
	ErrDatabase   = &Errno{Code: 20002, Message: "Database error."}
	ErrToken      = &Errno{Code: 20003, Message: "Error occurred while signing the json web token."}

	//user errors
	ErrEncrypt         = &Errno{Code: 20101, Message: "Error occurred while encrypt user password."}
	ErrTokenInvalid    = &Errno{Code: 20103, Message: "The token was invalid."}
	ErrPasswordInvalid = &Errno{Code: 20104, Message: "The password was invalid."}
)
