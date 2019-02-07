package auth

import (
	"github.com/lpimem/hlcsrv/storage"
)

type context_key int

const (
	AUTHENTICATED context_key = 1
	USER_ID       context_key = 2
	SESSION_ID    context_key = 3
	REASON        context_key = 4
)

const (
	HSESSION_ID = "x-hlc-token"
	HUSER_ID    = "x-hlc-uid"
)

// SessionInfo defines authentication session info
type SessionInfo struct {
	Uid storage.UserID `json:"uid"`
	Sid string         `json:"token"`
}
