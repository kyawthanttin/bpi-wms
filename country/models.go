package country

import (
	"errors"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/kyawthanttin/bpi-wms/dbutil"
	"github.com/kyawthanttin/bpi-wms/validation"
	validator "gopkg.in/go-playground/validator.v9"
)

type Country struct {
	Id   int    `json:"id" validate:"-"`
	Name string `json:"name" dbop:"iu" validate:"strmin=1,strmax=50,alphanumspecial"`
}

func ListCountries(db *sqlx.DB, search string) ([]Country, error) {
	results := []Country{}
	var err error

	if search != "" {
		s := Country{Name: search}
		nstmt, _ := db.PrepareNamed("SELECT id, name FROM Country WHERE UPPER(name) LIKE CONCAT('%', UPPER(:name), '%') ORDER BY name LIMIT " + strconv.Itoa(dbutil.MaxResults))
		err = nstmt.Select(&results, s)
	} else {
		err = db.Select(&results, "SELECT id,name FROM Country ORDER BY name LIMIT "+strconv.Itoa(dbutil.MaxResults))
	}
	return results, err
}

func GetCountry(db *sqlx.DB, id int) (Country, error) {
	result := Country{}
	err := db.Get(&result, "SELECT id, name FROM Country WHERE id = $1", id)
	return result, err
}

func CreateCountry(db *sqlx.DB, validate *validator.Validate, data Country) (Country, error) {
	if err := validate.Struct(data); err != nil {
		return Country{}, validation.DescribeErrors(err.(validator.ValidationErrors))
	}
	if exist, _ := dbutil.IsExist(db, "Country", "name", data.Name); exist {
		return Country{}, errors.New("Same country already exists")
	}
	id, err := dbutil.Insert(db, "Country", &data)
	if err != nil {
		return Country{}, err
	}
	return GetCountry(db, id.(int))
}

func UpdateCountry(db *sqlx.DB, validate *validator.Validate, id int, data Country) (Country, error) {
	if exist, _ := dbutil.IsExist(db, "Country", "id", id); !exist {
		return Country{}, errors.New("No such country")
	}
	if err := validate.Struct(data); err != nil {
		return Country{}, validation.DescribeErrors(err.(validator.ValidationErrors))
	}
	if exist, _ := dbutil.IsExistExcept(db, "Country", id, "name", data.Name); exist {
		return Country{}, errors.New("Same country already exists")
	}
	err := dbutil.Update(db, "Country", &data, &Country{Id: id})
	if err != nil {
		return Country{}, err
	}
	return GetCountry(db, id)
}

func DeleteCountry(db *sqlx.DB, id int) error {
	if exist, _ := dbutil.IsExist(db, "Country", "id", id); !exist {
		return errors.New("No such country")
	}
	_, err := db.Exec("DELETE FROM Country WHERE id = $1", id)
	return err
}
