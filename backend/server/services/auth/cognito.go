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

	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type AwsCognitorProvider struct {
	jwks     Jwks
	logger   log.Logger
	client   *cognitoidentityprovider.CognitoIdentityProvider
	clientId *string
}

func NewCognitoProvider(basicRes context.BasicRes) *AwsCognitorProvider {
	// Get configuration
	v := config.GetConfig()
	// TODO: verify the configuration
	// Create an AWS session
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(v.GetString("AWS_AUTH_REGION")),
	}))
	// Create a Cognito Identity Provider client
	client := cognitoidentityprovider.New(sess)
	cgt := &AwsCognitorProvider{
		client:   client,
		clientId: aws.String(v.GetString("AWS_AUTH_USER_POOL_WEB_CLIENT_ID")),
		logger:   basicRes.GetLogger().Nested("cognito"),
	}
	// Fetch the JWKS from the Cognito User Pool
	jwksURL := fmt.Sprintf(
		"https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json",
		v.GetString("AWS_AUTH_REGION"),
		v.GetString("AWS_AUTH_USER_POOL_ID"),
	)
	err := cgt.fetchJWKS(jwksURL)
	if err != nil {
		panic(err)
	}
	return cgt
}

func (cgt *AwsCognitorProvider) fetchJWKS(jwksURL string) errors.Error {
	// Get the JWKS from the URL
	resp, err := http.Get(jwksURL)
	if err != nil {
		return errors.Default.Wrap(err, "Failed to fetch JWKS")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Default.Wrap(err, "Failed to read JWKS")
	}
	// Unmarshal the response into a Jwks struct
	err = json.Unmarshal(body, &cgt.jwks)
	if err != nil {
		return errors.Default.Wrap(err, "Failed to unmarshall JWKS")
	}
	return nil
}

func (cgt *AwsCognitorProvider) SignIn(loginReq *LoginRequest) (*LoginResponse, errors.Error) {
	// Create the input for InitiateAuth
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		ClientId: cgt.clientId,
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(loginReq.Username),
			"PASSWORD": aws.String(loginReq.Password),
		},
	}

	// Call Cognito to get auth tokens
	response, err := cgt.client.InitiateAuth(input)
	if err != nil {
		return nil, errors.BadInput.New(err.Error())
	}

	loginRes := &LoginResponse{
		ChallengeName:       response.ChallengeName,
		ChallengeParameters: response.ChallengeParameters,
		Session:             response.Session,
	}
	if response.AuthenticationResult != nil {
		loginRes.AuthenticationResult = &AuthenticationResult{
			AccessToken:  response.AuthenticationResult.AccessToken,
			ExpiresIn:    response.AuthenticationResult.ExpiresIn,
			IdToken:      response.AuthenticationResult.IdToken,
			RefreshToken: response.AuthenticationResult.RefreshToken,
			TokenType:    response.AuthenticationResult.TokenType,
		}
	}

	return loginRes, nil
}

func (cgt *AwsCognitorProvider) CheckAuth(ctx *gin.Context) errors.Error {
	// Get the Auth header
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return errors.Unauthorized.New("Authorization header is missing")
	}

	// Split the header into "Bearer" and the actual token
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		return errors.Unauthorized.New("Invalid Authorization header")
	}

	// Parse the JWT token
	token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.Unauthorized.New(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
		}

		// Get the key ID from the header
		kid := token.Header["kid"].(string)

		// Look for the key that matches the kid
		for _, key := range cgt.jwks.Keys {
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
		cgt.logger.Error(err, "Invalid token")
		return errors.Unauthorized.New("Invalid token")
	}
	ctx.Set("token", token)
	return nil
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

func (cgt *AwsCognitorProvider) NewPassword(newPasswordReq *NewPasswordRequest) (*LoginResponse, errors.Error) {
	input := &cognitoidentityprovider.RespondToAuthChallengeInput{
		ChallengeName: aws.String("NEW_PASSWORD_REQUIRED"),
		ChallengeResponses: map[string]*string{
			"USERNAME":     aws.String(newPasswordReq.Username),
			"NEW_PASSWORD": aws.String(newPasswordReq.NewPassword),
		},
		Session:  aws.String(newPasswordReq.Session),
		ClientId: cgt.clientId,
	}
	response, err := cgt.client.RespondToAuthChallenge(input)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "Error setting up new password: "+err.Error())
	}
	// yes , it is identical to the login response, and yet they are 2 different structs
	loginRes := &LoginResponse{
		ChallengeName:       response.ChallengeName,
		ChallengeParameters: response.ChallengeParameters,
		Session:             response.Session,
	}
	if response.AuthenticationResult != nil {
		loginRes.AuthenticationResult = &AuthenticationResult{
			AccessToken:  response.AuthenticationResult.AccessToken,
			ExpiresIn:    response.AuthenticationResult.ExpiresIn,
			IdToken:      response.AuthenticationResult.IdToken,
			RefreshToken: response.AuthenticationResult.RefreshToken,
			TokenType:    response.AuthenticationResult.TokenType,
		}
	}
	return loginRes, nil
}

// func (cgt *AwsCognitorProvider) ChangePassword(ctx *gin.Context, oldPassword, newPassword string) errors.Error {
// 	token := ctx.GetString(("token"))
// 	if token == "" {
// 		return errors.Unauthorized.New("Token is missing")
// 	}
// 	input := &cognitoidentityprovider.ChangePasswordInput{
// 		AccessToken:      &token,
// 		PreviousPassword: &oldPassword,
// 		ProposedPassword: &newPassword,
// 	}
// 	_, err := cgt.client.ChangePassword(input)
// 	if err != nil {
// 		return errors.BadInput.Wrap(err, "Error changing password")
// 	}
// 	return nil
// }
