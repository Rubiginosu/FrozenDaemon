package auth

/**
保管了一个密钥对
*/
type KeyPair struct {
	ID   int   // 该ID对应服务器。
	Time int64 // 过期时间
}
