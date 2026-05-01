package core

import "codebase-app/pkg/errmsg"

func codedError(status int, message string, code string) *errmsg.CustomError {
	return errmsg.NewCustomErrors(status).
		SetMessage(message).
		SetErrorCode(code)
}
