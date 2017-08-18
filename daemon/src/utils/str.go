package utils

import "strings"

type CString string
func (s *CString)Contains(sep string) bool {
	return strings.Index(string(*s),sep) >= 0
}