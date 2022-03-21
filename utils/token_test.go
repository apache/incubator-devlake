package utils

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestGetEncodedToken(t *testing.T) {
	//v := config.GetConfig()
	fmt.Println(base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", "merico", "112db06547436eaa5b6d60da71a593ae8a"))))
}
