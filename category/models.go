package category

import (
	"errors"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/kyawthanttin/bpi-wms/dbutil"
)

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name" dbop:"iu"`
}

func ListCategories(db *sqlx.DB, search string) ([]Category, error) {
	results := []Category{}
	var err error

	if search != "" {
		s := Category{Name: search}
		nstmt, _ := db.PrepareNamed("SELECT id, name FROM Category WHERE UPPER(name) LIKE CONCAT('%', UPPER(:name), '%') ORDER BY name LIMIT " + strconv.Itoa(dbutil.MaxResults))
		err = nstmt.Select(&results, s)
	} else {
		err = db.Select(&results, "SELECT id,name FROM Category ORDER BY name LIMIT "+strconv.Itoa(dbutil.MaxResults))
	}
	return results, err
}

func GetCategory(db *sqlx.DB, id int) (Category, error) {
	result := Category{}
	err := db.Get(&result, "SELECT id, name FROM Category WHERE id = $1", id)
	return result, err
}

func CreateCategory(db *sqlx.DB, data Category) (Category, error) {
	if exist, _ := dbutil.IsExist(db, "Category", "name", data.Name); exist {
		return Category{}, errors.New("Same category already exists")
	}
	id, err := dbutil.Insert(db, "Category", &data)
	if err != nil {
		return Category{}, err
	}
	return GetCategory(db, id.(int))
}

func UpdateCategory(db *sqlx.DB, id int, data Category) (Category, error) {
	if exist, _ := dbutil.IsExist(db, "Category", "id", id); !exist {
		return Category{}, errors.New("No such category")
	}
	err := dbutil.Update(db, "Category", &data, &Category{Id: id})
	if err != nil {
		return Category{}, err
	}
	return GetCategory(db, id)
}

func DeleteCategory(db *sqlx.DB, id int) error {
	if exist, _ := dbutil.IsExist(db, "Category", "id", id); !exist {
		return errors.New("No such category")
	}
	_, err := db.Exec("DELETE FROM Category WHERE id = $1", id)
	return err
}
