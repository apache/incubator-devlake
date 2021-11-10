package models

type AuthToken struct {
	Model
	UserName string
	Token    string
}

func GetAuthToken() ([]AuthToken, error){
	var tokens []AuthToken
	err := Db.Find(&tokens).Error
	return tokens, err
}