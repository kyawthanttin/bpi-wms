package user

import (
	"errors"
	"strconv"
	"time"

	"github.com/leebenson/conform"

	"github.com/jmoiron/sqlx"
	"github.com/kyawthanttin/bpi-wms/dbutil"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int       `json:"id"`
	Username  string    `json:"username" conform:"snake" dbop:"iu"`
	Password  string    `json:"password" dbop:"i"`
	Name      string    `json:"name" conform:"name" dbop:"iu"`
	Roles     string    `json:"roles" dbop:"iu"`
	IsEnabled bool      `json:"isEnabled" db:"is_enabled" dbop:"iu"`
	Created   time.Time `json:"created" dbop:"i"`
	LastLogin time.Time `json:"lastLogin" db:"last_login" dbop:"i"`
}

func ListUsers(db *sqlx.DB, search string) ([]User, error) {
	results := []User{}
	var err error

	if search != "" {
		s := User{Username: search, Name: search, Roles: search}
		nstmt, _ := db.PrepareNamed("SELECT id, username, name, roles, is_enabled, last_login FROM Login_User WHERE UPPER(username) LIKE CONCAT('%', UPPER(:username), '%') " +
			" OR UPPER(name) LIKE CONCAT('%', UPPER(:name), '%') OR UPPER(roles) LIKE CONCAT('%', UPPER(:roles), '%') ORDER BY name LIMIT " + strconv.Itoa(dbutil.MaxResults))
		err = nstmt.Select(&results, s)
	} else {
		err = db.Select(&results, "SELECT id, username, name, roles, is_enabled, last_login FROM Login_User ORDER BY name LIMIT "+strconv.Itoa(dbutil.MaxResults))
	}
	return results, err
}

func GetUser(db *sqlx.DB, id int) (User, error) {
	result := User{}
	err := db.Get(&result, "SELECT id, username, name, roles, is_enabled, created, last_login FROM Login_User WHERE id = $1", id)
	return result, err
}

func GetUserByUsername(db *sqlx.DB, username string) (User, error) {
	result := User{}
	err := db.Get(&result, "SELECT id, username, password, name, roles, is_enabled, created, last_login FROM Login_User WHERE username = $1", username)
	return result, err
}

func CreateUser(db *sqlx.DB, data User) (User, error) {
	conform.Strings(&data)
	if exist, _ := dbutil.IsExist(db, "Login_User", "username", data.Username); exist {
		return User{}, errors.New("Same username already exists")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.MinCost)
	if err != nil {
		return User{}, err
	}
	data.Password = string(hashPassword)
	data.Created = time.Now()
	data.LastLogin = time.Now()
	id, err := dbutil.Insert(db, "Login_User", &data)
	if err != nil {
		return User{}, err
	}
	return GetUser(db, id.(int))
}

func UpdateUser(db *sqlx.DB, id int, data User) (User, error) {
	if exist, _ := dbutil.IsExist(db, "Login_User", "id", id); !exist {
		return User{}, errors.New("No such user")
	}
	conform.Strings(&data)
	err := dbutil.Update(db, "Login_User", &data, &User{Id: id})
	if err != nil {
		return User{}, err
	}
	return GetUser(db, id)
}

func DeleteUser(db *sqlx.DB, id int) error {
	if exist, _ := dbutil.IsExist(db, "Login_User", "id", id); !exist {
		return errors.New("No such user")
	}
	_, err := db.Exec("DELETE FROM Login_User WHERE id = $1", id)
	return err
}

func ChangePassword(db *sqlx.DB, id int, password string) (User, error) {
	if exist, _ := dbutil.IsExist(db, "Login_User", "id", id); !exist {
		return User{}, errors.New("No such user")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return User{}, err
	}
	_, err = db.Exec("UPDATE Login_User SET password = $1 WHERE id = $2", string(hashPassword), id)
	if err != nil {
		return User{}, err
	}
	return GetUser(db, id)
}

func RecordLogin(db *sqlx.DB, id int) error {
	_, err := db.Exec("UPDATE Login_User SET last_login = $1 WHERE id = $2", time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}
