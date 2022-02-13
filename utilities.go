package main

import (
	"os"
	"io/ioutil"
	"fmt"
)

// safeWrite throws errors if the file already exists
func safeWrite(path string, data []byte, perms os.FileMode, override bool) error {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return ioutil.WriteFile(path, data, perms)
	}

	if override == true {
		return ioutil.WriteFile(path, data, perms)
	}

	return fmt.Errorf("cannot write \"%s\"; file already exists", path)
}
