package supplier

import (
	"database/sql"
)

type Supplier struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

func ListSuppliers(db *sql.DB, search string) ([]*Supplier, error) {
	var rows *sql.Rows
	var err error

	if search != "" {
		rows, err = db.Query("SELECT id, name, address FROM Supplier WHERE UPPER(name) LIKE CONCAT('%', UPPER($1), '%') ORDER BY name", search)
	} else {
		rows, err = db.Query("SELECT id, name, address FROM Supplier ORDER BY name")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*Supplier, 0)
	for rows.Next() {
		result := new(Supplier)
		err := rows.Scan(&result.Id, &result.Name, &result.Address)
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

func GetSupplier(db *sql.DB, id int) (*Supplier, error) {
	row := db.QueryRow("SELECT id, name, address FROM Supplier WHERE id = $1", id)
	result := new(Supplier)
	err := row.Scan(&result.Id, &result.Name, &result.Address)
	if err != nil {
		return &Supplier{}, err
	}
	return result, nil
}

func CreateSupplier(db *sql.DB, data Supplier) (*Supplier, error) {
	sqlStatement := `INSERT INTO Supplier(name, address) VALUES ($1, $2) RETURNING id`
	var id int
	err := db.QueryRow(sqlStatement, data.Name, data.Address).Scan(&id)
	if err != nil {
		return &Supplier{}, err
	}
	return GetSupplier(db, id)
}

func UpdateSupplier(db *sql.DB, id int, data Supplier) (*Supplier, error) {
	_, err := db.Exec("UPDATE Supplier SET name = $1, address = $2 WHERE id = $3", data.Name, data.Address, id)
	if err != nil {
		return &Supplier{}, err
	}
	return GetSupplier(db, id)
}

func DeleteSupplier(db *sql.DB, id int) error {
	if _, err := GetSupplier(db, id); err != nil {
		return err
	}

	result, err := db.Exec("DELETE FROM Supplier WHERE id = $1", id)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}
