package ehelper

import "os"

// Check whether the file exists
func FileIsExist(path string) bool {
	_, err := os.Stat(path)

	return err == nil || os.IsExist(err)
}
