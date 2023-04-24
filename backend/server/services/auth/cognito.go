/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync"

	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/impls/logruslog"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	//jwksCache is a cache of the fetched JWKS
	jwksCache Jwks
	//jwksCacheMx is a mutex to lock the jwksCache
	jwksCacheMx sync.Mutex
	logger      = logruslog.Global.Nested("auth")
)

func CreateCognitoClient() *cognitoidentityprovider.CognitoIdentityProvider {
	// Get configuration
	v := config.GetConfig()
	// Create an AWS session
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(v.GetString("AWS_AUTH_REGION")),
	}))
	// Create a Cognito Identity Provider client
	return cognitoidentityprovider.New(sess)
}

func SignIn(cognitoClient *cognitoidentityprovider.CognitoIdentityProvider, username, password string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	// Get configuration
	v := config.GetConfig()
	// Create the input for InitiateAuth
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		ClientId: aws.String(v.GetString("AWS_AUTH_USER_POOL_WEB_CLIENT_ID")),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(username),
			"PASSWORD": aws.String(password),
		},
	}

	// Call Cognito to get auth tokens
	response, err := cognitoClient.InitiateAuth(input)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func fetchJWKS(jwksURL string) (jwks Jwks, err error) {
	// Get the JWKS from the URL
	resp, err := http.Get(jwksURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	// Unmarshal the response into a Jwks struct
	err = json.Unmarshal(body, &jwks)
	return
}

func ensureJWKS(jwksURL string) (jwks Jwks, err error) {
	// Lock the mutex
	jwksCacheMx.Lock()
	defer jwksCacheMx.Unlock()

	// If the cache is empty, fetch the JWKS
	if len(jwksCache.Keys) == 0 {
		jwksCache, err = fetchJWKS(jwksURL)
	}
	// Return the cached JWKS
	jwks = jwksCache
	return
}

func AuthenticationMiddleware(ctx *gin.Context) {
	// Get configuration
	v := config.GetConfig()
	// Construct the JWKS URL
	jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", v.GetString("AWS_AUTH_REGION"), v.GetString("AWS_AUTH_USER_POOL_ID"))
	// Get the cached JWKS
	jwks, err := ensureJWKS(jwksURL)
	if err != nil {
		fmt.Printf("Error fetching JWKS: %v\n", err)
		ctx.Abort()
		return
	}

	// Get the Auth header
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		http.Error(ctx.Writer, "Authorization header is missing", http.StatusUnauthorized)
		ctx.Abort()
		return
	}

	// Split the header into "Bearer" and the actual token
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		http.Error(ctx.Writer, "Invalid Authorization header", http.StatusUnauthorized)
		ctx.Abort()
		return
	}

	// Parse the JWT token
	token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Get the key ID from the header
		kid := token.Header["kid"].(string)

		// Look for the key that matches the kid
		for _, key := range jwks.Keys {
			if key.Kid == kid {
				// Construct the RSA public key
				n := pemHeader(key.N)
				e := pemHeader(key.E)
				parsedKey := &rsa.PublicKey{
					N: new(big.Int).SetBytes(n),
					E: int(new(big.Int).SetBytes(e).Int64()),
				}
				return parsedKey, nil
			}
		}

		return nil, fmt.Errorf("Public key not found")
	})

	// Check if the token is invalid
	if err != nil || !token.Valid {
		logger.Error(err, "Invalid token")
		http.Error(ctx.Writer, "Invalid token", http.StatusUnauthorized)
		ctx.Abort()
	}
}

func pemHeader(encodedKey string) []byte {
	// Decode the base64 encoded key
	key, err := base64.RawURLEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil
	}
	return key
}

// Jwks represents the JSON web key set
type Jwks struct {
	Keys []struct {
		Kid string `json:"kid"`
		N   string `json:"n"`
		E   string `json:"e"`
	} `json:"keys"`
}
