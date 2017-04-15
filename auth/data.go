package auth

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

// defines authentication session info
type SessionInfo struct {
	Uid uint32 `json:"uid"`
	Sid string `json:"token"`
}
