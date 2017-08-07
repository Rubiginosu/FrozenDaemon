package auth

import "time"

var keys = make(map[string]*KeyPair, 0)

const (
	KEY_OUT_DATE = -2 + iota
	KEY_VERIFY_FAILED
)

// 传入Key,返回服务器id或一些状态值
func VerifyKey(key string) int {
	if pair, ok := keys[key]; ok {
		if pair.Time <= time.Now().Unix() {
			delete(keys, key)
			return KEY_OUT_DATE
		}
		return pair.ID
	}
	return KEY_VERIFY_FAILED
}

func KeyRigist(key string, id int) {
	keys[key] = &KeyPair{
		ID:   id,
		Time: time.Now().Unix() + 600,
	}
}
