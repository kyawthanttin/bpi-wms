package supplier

import (
	"errors"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/kyawthanttin/bpi-wms/dbutil"
)

type Supplier struct {
	Id      int    `json:"id"`
	Name    string `json:"name" dbop:"iu"`
	Address string `json:"address" dbop:"iu"`
}

func ListSuppliers(db *sqlx.DB, search string) ([]Supplier, error) {
	results := []Supplier{}
	var err error

	if search != "" {
		s := Supplier{Name: search}
		nstmt, _ := db.PrepareNamed("SELECT id, name, address FROM Supplier WHERE UPPER(name) LIKE CONCAT('%', UPPER(:name), '%') ORDER BY name LIMIT " + strconv.Itoa(dbutil.MaxResults))
		err = nstmt.Select(&results, s)
	} else {
		err = db.Select(&results, "SELECT id, name, address FROM Supplier ORDER BY name LIMIT "+strconv.Itoa(dbutil.MaxResults))
	}
	return results, err
}

func GetSupplier(db *sqlx.DB, id int) (Supplier, error) {
	result := Supplier{}
	err := db.Get(&result, "SELECT id, name, address FROM Supplier WHERE id = $1", id)
	return result, err
}

func CreateSupplier(db *sqlx.DB, data Supplier) (Supplier, error) {
	if exist, _ := dbutil.IsExist(db, "Supplier", "name", data.Name); exist {
		return Supplier{}, errors.New("Same supplier already exists")
	}
	id, err := dbutil.Insert(db, "Supplier", &data)
	if err != nil {
		return Supplier{}, err
	}
	return GetSupplier(db, id.(int))
}

func UpdateSupplier(db *sqlx.DB, id int, data Supplier) (Supplier, error) {
	if exist, _ := dbutil.IsExist(db, "Supplier", "id", id); !exist {
		return Supplier{}, errors.New("No such supplier")
	}
	err := dbutil.Update(db, "Supplier", &data, &Supplier{Id: id})
	if err != nil {
		return Supplier{}, err
	}
	return GetSupplier(db, id)
}

func DeleteSupplier(db *sqlx.DB, id int) error {
	if exist, _ := dbutil.IsExist(db, "Supplier", "id", id); !exist {
		return errors.New("No such supplier")
	}
	_, err := db.Exec("DELETE FROM Supplier WHERE id = $1", id)
	return err
}
