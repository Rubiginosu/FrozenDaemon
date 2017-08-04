package auth

func UserAuth(userServerID int, dst string, index int) bool {
	if ValidationKeyPairs[index].ValidationKeyPair.ID != userServerID {
		return false
	}
	return ValidationKeyPairs[index].ValidationKeyPair.Key == dst
}
