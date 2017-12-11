package authentication

import (
	"database/sql"
)

type UserInfo struct {
	Id       int
	Username string `json:"username"`
	Name     string `json:"name"`
}

func ValidateUser(db *sql.DB, username string, password string) (*UserInfo, error) {
	row := db.QueryRow("SELECT id, username, name FROM login_user WHERE username = $1 AND password = $2", username, password)
	userInfo := new(UserInfo)
	err := row.Scan(&userInfo.Id, &userInfo.Username, &userInfo.Name)
	if err != nil {
		return &UserInfo{}, err
	}
	return userInfo, nil
}
