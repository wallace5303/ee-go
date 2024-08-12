package ehelper

import (
	"math/rand"
	"time"
)

// String returns a random string ['a', 'z'] and ['0', '9'] in the specified length.
func GetRandomString(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	time.Sleep(10 * time.Nanosecond)

	letter := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

// Contains determines whether the str is in the strs.
func Contains(str string, strs []string) bool {
	for _, v := range strs {
		if v == str {
			return true
		}
	}
	return false
}

// RemoveDuplicatedElem removes the duplicated elements from the slice.
func RemoveDuplicatedElem(slice []string) (ret []string) {
	allKeys := make(map[string]bool)
	for _, item := range slice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			ret = append(ret, item)
		}
	}
	return
}
