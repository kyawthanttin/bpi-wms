package authentication

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	Id        int       `json:"id" validate:"-"`
	Username  string    `json:"username" dbop:"iu" validate:"username"`
	Password  string    `json:"password" dbop:"i" validate:"omitempty,password"`
	Name      string    `json:"name" dbop:"iu" validate:"strmin=1,strmax=50,alphanumspecial"`
	Roles     string    `json:"roles" dbop:"iu" validate:"strmin=1"`
	IsEnabled bool      `json:"isEnabled" db:"is_enabled" dbop:"iu" validate:"-"`
	Created   time.Time `json:"created" dbop:"i" validate:"-"`
	LastLogin time.Time `json:"lastLogin" db:"last_login" dbop:"i" validate:"-"`
}

func GetUserByUsername(db *sqlx.DB, username string) (User, error) {
	result := User{}
	err := db.Get(&result, "SELECT id, username, password, name, roles, is_enabled, created, last_login FROM Login_User WHERE username = $1", username)
	return result, err
}

func RecordLogin(db *sqlx.DB, id int) error {
	_, err := db.Exec("UPDATE Login_User SET last_login = $1 WHERE id = $2", time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}
