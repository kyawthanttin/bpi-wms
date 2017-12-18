package item

import (
	"errors"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kyawthanttin/bpi-wms/dbutil"
	"github.com/kyawthanttin/bpi-wms/validation"
	validator "gopkg.in/go-playground/validator.v9"
)

type Item struct {
	Id           int       `json:"id" validate:"-"`
	Code         string    `json:"code" dbop:"i" validate:"strmin=1,strmax=12,alphanum"`
	Name         string    `json:"name" dbop:"iu" validate:"strmin=1,strmax=100,alphanumspecial"`
	Type         string    `json:"type" dbop:"iu" validate:"strmin=1,strmax=100,alphanumspecial"`
	BrandName    string    `json:"brandName" db:"brand_name" dbop:"iu" validate:"strmax=100,alphanumspecial"`
	PackingSize  string    `json:"packingSize" db:"packing_size" dbop:"iu" validate:"strmax=100,alphanumspecial"`
	CasePack     float32   `json:"casePack" db:"case_pack" dbop:"iu" validate:"required,gt=0"`
	CaseUom      int       `json:"caseUom" db:"caseuom" dbop:"iu" validate:"required,gt=0"`
	CaseUomName  string    `json:"caseUomName" db:"caseuomname" validate:"-"`
	PieceUom     int       `json:"pieceUom" db:"pieceuom" dbop:"iu" validate:"required,gt=0"`
	PieceUomName string    `json:"pieceUomName" db:"pieceuomname" validate:"-"`
	Category     int       `json:"category" dbop:"iu" validate:"required,gt=0"`
	CategoryName string    `json:"categoryName" db:"categoryname" validate:"-"`
	Created      time.Time `json:"created" dbop:"i" validate:"-"`
	LastModified time.Time `json:"lastModified" db:"last_modified" dbop:"iu" validate:"-"`
}

func ListItems(db *sqlx.DB, search string) ([]Item, error) {
	results := []Item{}
	var err error

	if search != "" {
		s := Item{Code: search, Name: search}
		nstmt, _ := db.PrepareNamed("SELECT id, code, name, type, brand_name, packing_size, case_pack, caseuom, pieceuom, category FROM Item " +
			"WHERE UPPER(code) LIKE CONCAT('%', UPPER(:code), '%') OR UPPER(name) LIKE CONCAT('%', UPPER(:name), '%') ORDER BY name LIMIT " + strconv.Itoa(dbutil.MaxResults))
		err = nstmt.Select(&results, s)
	} else {
		err = db.Select(&results, "SELECT id, code, name, type, brand_name, packing_size, case_pack, caseuom, pieceuom, category FROM Item ORDER BY name LIMIT "+
			strconv.Itoa(dbutil.MaxResults))
	}
	return results, err
}

func GetItem(db *sqlx.DB, id int) (Item, error) {
	result := Item{}
	err := db.Get(&result, "SELECT i.id, i.code, i.name, i.type, i.brand_name, i.packing_size, i.case_pack, i.caseuom, i.pieceuom, i.category, i.created, i.last_modified, "+
		"c.name AS categoryname, u1.name || ' (' || u1.abbr || ')' AS caseuomname, u2.name || ' (' || u2.abbr || ')' AS pieceuomname "+
		"FROM Item i, Category c, UnitOfMeasurement u1, UnitOfMeasurement u2 "+
		"WHERE i.category = c.id AND i.caseuom = u1.id AND i.pieceuom = u2.id AND i.id = $1", id)
	return result, err
}

func CreateItem(db *sqlx.DB, validate *validator.Validate, data Item) (Item, error) {
	if err := validate.Struct(data); err != nil {
		return Item{}, validation.DescribeErrors(err.(validator.ValidationErrors))
	}
	if exist, _ := dbutil.IsExist(db, "Item", "code", data.Code); exist {
		return Item{}, errors.New("Same item code already exists")
	}
	data.Created = time.Now()
	data.LastModified = time.Now()
	id, err := dbutil.Insert(db, "Item", &data)
	if err != nil {
		return Item{}, err
	}
	return GetItem(db, id.(int))
}

func UpdateItem(db *sqlx.DB, validate *validator.Validate, id int, data Item) (Item, error) {
	if exist, _ := dbutil.IsExist(db, "Item", "id", id); !exist {
		return Item{}, errors.New("No such item")
	}
	if err := validate.Struct(data); err != nil {
		return Item{}, validation.DescribeErrors(err.(validator.ValidationErrors))
	}
	data.LastModified = time.Now()
	err := dbutil.Update(db, "Item", &data, &Item{Id: id})
	if err != nil {
		return Item{}, err
	}
	return GetItem(db, id)
}

func DeleteItem(db *sqlx.DB, id int) error {
	if exist, _ := dbutil.IsExist(db, "Item", "id", id); !exist {
		return errors.New("No such item")
	}
	_, err := db.Exec("DELETE FROM Item WHERE id = $1", id)
	return err
}
