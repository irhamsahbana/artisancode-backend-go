package common

import (
	"codebase-app/pkg/errmsg"
	"errors"
)

var (
	ErrProductNotFound = errors.New(errmsg.MessageProductNotFound)
)

var (
	ErrCustomProductNotFound = errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageProductNotFound)
)
