package core

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"github.com/merico-dev/lake/config"
)

type TestResult struct {
	Success bool
	Message string
}

func (testResult *TestResult) Set(success bool, message string) {
	testResult.Success = success
	testResult.Message = message
}

func ValidateParams(input *ApiResourceInput, requiredParams []string) *TestResult {
	message := "Missing params: "
	missingParams := []string{}
	if len(input.Body) == 0 {
		for _, param := range requiredParams {
			message += fmt.Sprintf(" %v", param)
		}
		return &TestResult{Success: false, Message: message}
	} else {
		for _, param := range requiredParams {
			if input.Body[param] == "" {
				missingParams = append(missingParams, param)
			}
		}
		if len(missingParams) > 0 {
			for _, param := range missingParams {
				message += fmt.Sprintf(" %v", param)
			}
			return &TestResult{Success: false, Message: message}
		} else {
			return &TestResult{Success: true, Message: ""}
		}
	}
}

const InvalidParams = "Failed to decode request params"
const SourceIdError = "Missing or Invalid sourceId"
const InvalidConnectionError = "Your connection configuration is invalid."
const UnsetConnectionError = "Your connection configuration is not set."
const UnmarshallingError = "There was a problem unmarshalling the response"
const InvalidEndpointError = "Failed to parse endpoint"
const SchemaIsRequired = "Endpoint schema is required"
const InvalidSchema = "Failed to find port for schema"
const DNSResolveFailedError = "Failed to find ip address"
const NetworkConnectError = "Failed to connect to endpoint"
const EncodeKeyEnvStr = "ENCODE_KEY"

func GetRateLimitPerSecond(options map[string]interface{}, defaultValue int) (int, error) {
	if options["rateLimitPerSecond"] == nil {
		return defaultValue, nil
	}

	rateLimitPerSecond := options["rateLimitPerSecond"]
	if value, ok := rateLimitPerSecond.(float64); ok {
		return int(value), nil
	} else {
		return 0, fmt.Errorf("rateLimitPerSecond is invalid")
	}
}

// AES + Base64 encryption using ENCODE_KEY in .env as key
func Encode(Input string) (string, error) {
	// Read encryption key from configuration
	V := config.LoadConfigFile()
	encodingKey := V.GetString(EncodeKeyEnvStr)
	// when encryption key is not set
	if encodingKey == "" {
		// Randomly generate a bunch of encryption keys and set them to config
		encodingKey = RandomCapsStr(128)
		V.Set(EncodeKeyEnvStr, encodingKey)
		V.WriteConfig()
	}
	// add suffix to the data part
	inputBytes := append([]byte(Input), 123, 110, 100, 100, 116, 102, 125)
	// perform encryption
	output, err := AesEncrypt(inputBytes, []byte(encodingKey))
	if err != nil {
		return Input, err
	}
	// Return the result after Base64 processing
	return base64.StdEncoding.EncodeToString(output), nil
}

//  Base64 + AES decryption using ENCODE_KEY in .env as key
func Decode(Input string) (string, error) {
	// Read encryption key from configuration
	V := config.LoadConfigFile()
	encodingKey := V.GetString(EncodeKeyEnvStr)
	// when encryption key is not set
	if encodingKey == "" {
		// return error message
		return Input, fmt.Errorf("The setting ENCODE_KEY from the file '.env' is empty.decrypted fail.")
	}

	// Decode Base64
	decodingFromBase64, err1 := base64.StdEncoding.DecodeString(Input)
	if err1 != nil {
		return Input, err1
	}
	// perform AES decryption
	output, err2 := AesDecrypt(decodingFromBase64, []byte(encodingKey))
	if err2 != nil {
		return Input, err2
	}

	// Verify and remove suffix
	oSize := len(output)
	if oSize >= 7 {
		check := output[oSize-7 : oSize]
		backEnd := []byte{123, 110, 100, 100, 116, 102, 125}
		if string(check) == string(backEnd) {
			output = output[0 : oSize-7]
			// return result
			return string(output), nil
		}
	}
	return "", fmt.Errorf("The setting ENCODE_KEY from the file '.env' is incorrect.decrypted fail.")
}

// PKCS7 padding
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7 unPadding
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	if unpadding >= length {
		return nil
	}
	return origData[:(length - unpadding)]
}

//AES encryption, CBC
func AesEncrypt(origData, key []byte) ([]byte, error) {
	// data alignment fill and encryption
	sha256Key := sha256.Sum256(key)
	key = sha256Key[:]
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// data alignment fill and encryption
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//AES decryption
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	// Uniformly use sha256 to process as 32-bit Byte (256-bit bit)
	sha256Key := sha256.Sum256(key)
	key = sha256Key[:]
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// Decrypt and unalign data
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

// A random string of length len uppercase characters
func RandomCapsStr(len int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
