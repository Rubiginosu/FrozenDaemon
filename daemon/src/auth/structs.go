package auth

import "time"

type ValidationKeyPairTime struct {
	ValidationKeyPair ValidationKeyPair
	GeneratedTime     time.Time
}
type ValidationKeyPair struct {
	ID  int // 该ID对应服务器。
	Key string
}
