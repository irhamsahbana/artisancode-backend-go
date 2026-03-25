package common

import (
	"codebase-app/pkg/errmsg"
	"errors"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

var (
	ErrCustomProductNotFound = errmsg.NewCustomErrors(404).SetMessage("product not found")
)
