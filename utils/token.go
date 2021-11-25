package utils

import (
	"encoding/base64"
	"fmt"
)

func GetEncodedToken(username string, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", username, password)))
}
