package auth

import "time"

var keys = make(map[string]*KeyPair, 0)

const (
	KEY_VERIFY_FAILED = -1
)
func Timer(){
	for {
		for k,v := range keys {
			if v.Time <= time.Now().Unix() {
				delete(keys,k)
			}
		}
		time.Sleep(10 * time.Second)
	}
}
// 传入Key,返回服务器id或一些状态值
func VerifyKey(key string) int {
	if pair, ok := keys[key]; ok {
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
