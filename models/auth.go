package models

type AuthToken struct {
	Model
	UserName string
	Token    string
}

func GetAuthTokenList() ([]AuthToken, error){
	var tokens []AuthToken
	err := Db.Find(&tokens).Error
	return tokens, err
}

func GetAuthToken(token string) (*AuthToken, error){
	authToken := new(AuthToken)
	err := Db.First(authToken, "token = ?", token).Error
	return authToken, err
}