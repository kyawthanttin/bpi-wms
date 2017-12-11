package item

import (
	"database/sql"
	"time"
)

type Item struct {
	Id           int       `json:"id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	PackingSize  string    `json:"packingSize"`
	CasePack     float32   `json:"casePack"`
	CaseUomId    int       `json:"caseUomId"`
	CaseUom      string    `json:"caseUom"`
	PieceUomId   int       `json:"pieceUomId"`
	PieceUom     string    `json:"pieceUom"`
	CategoryId   int       `json:"categoryId"`
	Category     string    `json:"category"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
}

func ListItems(db *sql.DB, search string) ([]*Item, error) {
	var rows *sql.Rows
	var err error

	if search != "" {
		rows, err = db.Query("SELECT id, code, name, type, packing_size, case_pack, case_uom, piece_uom, category FROM Item "+
			"WHERE UPPER(code) LIKE CONCAT('%', UPPER($1), '%') OR UPPER(name) LIKE CONCAT('%', UPPER($1), '%') ORDER BY code", search)
	} else {
		rows, err = db.Query("SELECT id, code, name, type, packing_size, case_pack, case_uom, piece_uom, category FROM Item ORDER BY code")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*Item, 0)
	for rows.Next() {
		result := new(Item)
		err := rows.Scan(&result.Id, &result.Code, &result.Name, &result.Type, &result.PackingSize, &result.CasePack, &result.CaseUomId, &result.PieceUomId, &result.CategoryId)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func GetItem(db *sql.DB, id int) (*Item, error) {
	row := db.QueryRow("SELECT id, code, name, type, packing_size, case_pack, case_uom, piece_uom, category, created, last_modified FROM Item WHERE id = $1", id)
	result := new(Item)
	err := row.Scan(&result.Id, &result.Code, &result.Name, &result.Type, &result.PackingSize, &result.CasePack, &result.CaseUomId, &result.PieceUomId, &result.CategoryId, &result.Created, &result.LastModified)
	if err != nil {
		return &Item{}, err
	}
	return result, nil
}

func CreateItem(db *sql.DB, data Item) (*Item, error) {
	sqlStatement := `INSERT INTO Item(code, name, type, packing_size, case_pack, case_uom, piece_uom, category, created, last_modified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`
	var id int
	err := db.QueryRow(sqlStatement, data.Code, data.Name, data.Type, data.PackingSize, data.CasePack, data.CaseUomId, data.PieceUomId, data.CategoryId, time.Now(), time.Now()).Scan(&id)
	if err != nil {
		return &Item{}, err
	}
	return GetItem(db, id)
}

func UpdateItem(db *sql.DB, id int, data Item) (*Item, error) {
	_, err := db.Exec("UPDATE Item SET code = $1, name = $2, type = $3, packing_size = $4, case_pack = $5, case_uom = $6, piece_uom = $7, category = $8, last_modified = $9 WHERE id = $10",
		data.Code, data.Name, data.Type, data.PackingSize, data.CasePack, data.CaseUomId, data.PieceUomId, data.CategoryId, time.Now(), id)
	if err != nil {
		return &Item{}, err
	}
	return GetItem(db, id)
}

func DeleteItem(db *sql.DB, id int) error {
	if _, err := GetItem(db, id); err != nil {
		return err
	}

	result, err := db.Exec("DELETE FROM Item WHERE id = $1", id)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}
