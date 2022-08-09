package config

import (
	"io/ioutil"
)

//LoadKeyFromFile function, load key from file and parse into []byte
func LoadKeyFromFile(filePath string) ([]byte, error) {
	signBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return signBytes, nil
}
