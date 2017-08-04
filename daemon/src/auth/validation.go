package auth

import (
	"conf"
	"time"
)

var ValidationKeyPairs []ValidationKeyPairTime

func ValidationKeyGenerate(id int) ValidationKeyPairTime {
	pair := ValidationKeyPairTime{
		ValidationKeyPair: ValidationKeyPair{
			ID:  id,
			Key: conf.RandString(20),
		},
		GeneratedTime: time.Now(),
	}
	return pair
}
func ValidationKeyUpdate(outDateSeconds float64) {
	for {
		validationKeyClear(outDateSeconds)
		time.Sleep(300 * time.Second)
	}
}
func validationKeyClear(outDateSeconds float64) {
	j := 0
	i := 0
	for k := j; k < len(ValidationKeyPairs); k++ {
		if isValidationKeyAvailable(ValidationKeyPairs[k], outDateSeconds) {
			// swap [swapper] and [k]
			temp := ValidationKeyPairs[i]
			ValidationKeyPairs[i] = ValidationKeyPairs[k]
			ValidationKeyPairs[k] = temp
			// i指针自增
			i++
		}
	}
	ValidationKeyPairs = ValidationKeyPairs[i:]
}

func isValidationKeyAvailable(pairs ValidationKeyPairTime, outDateSeconds float64) bool {
	return time.Since(pairs.GeneratedTime).Seconds() > outDateSeconds
}

func FindValidationKey(target int) int {
	for i := 0; i < len(ValidationKeyPairs); i++ {
		if ValidationKeyPairs[i].ValidationKeyPair.ID == target {
			return i
		}
	}
	return -1
}

func GetValidationKeyPairs() []ValidationKeyPairTime {
	return ValidationKeyPairs
}
func IsVerifiedValidationKeyPair(id int, key string) bool {
	if i := FindValidationKey(id); i > -1 {
		return ValidationKeyPairs[i].ValidationKeyPair.Key == key
	}
	return false
}
