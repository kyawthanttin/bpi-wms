package country

import (
	"errors"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/kyawthanttin/bpi-wms/dbutil"
)

type Country struct {
	Id   int    `json:"id"`
	Name string `json:"name" dbop:"iu"`
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

func CreateCountry(db *sqlx.DB, data Country) (Country, error) {
	if exist, _ := dbutil.IsExist(db, "Country", "name", data.Name); exist {
		return Country{}, errors.New("Same country already exists")
	}
	id, err := dbutil.Insert(db, "Country", &data)
	if err != nil {
		return Country{}, err
	}
	return GetCountry(db, id.(int))
}

func UpdateCountry(db *sqlx.DB, id int, data Country) (Country, error) {
	if exist, _ := dbutil.IsExist(db, "Country", "id", id); !exist {
		return Country{}, errors.New("No such country")
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
