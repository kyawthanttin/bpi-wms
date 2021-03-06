package user

import (
	"errors"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kyawthanttin/bpi-wms/dbutil"
	"github.com/kyawthanttin/bpi-wms/validation"
	validator "gopkg.in/go-playground/validator.v9"

	"golang.org/x/crypto/bcrypt"
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

func CreateUser(db *sqlx.DB, validate *validator.Validate, data User) (User, error) {
	if err := validate.Struct(data); err != nil {
		return User{}, validation.DescribeErrors(err.(validator.ValidationErrors))
	}
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

func UpdateUser(db *sqlx.DB, validate *validator.Validate, id int, data User) (User, error) {
	if exist, _ := dbutil.IsExist(db, "Login_User", "id", id); !exist {
		return User{}, errors.New("No such user")
	}
	if err := validate.Struct(data); err != nil {
		return User{}, validation.DescribeErrors(err.(validator.ValidationErrors))
	}
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

func ChangePassword(db *sqlx.DB, validate *validator.Validate, id int, password string) (User, error) {
	if exist, _ := dbutil.IsExist(db, "Login_User", "id", id); !exist {
		return User{}, errors.New("No such user")
	}
	if err := validate.Var(password, "password"); err != nil {
		return User{}, validation.DescribeErrors(err.(validator.ValidationErrors))
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
