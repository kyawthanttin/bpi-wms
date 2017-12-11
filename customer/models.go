package customer

import (
	"database/sql"
)

type Customer struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func ListCustomers(db *sql.DB, search string) ([]*Customer, error) {
	var rows *sql.Rows
	var err error

	if search != "" {
		rows, err = db.Query("SELECT id, name FROM Customer WHERE UPPER(name) LIKE CONCAT('%', UPPER($1), '%') ORDER BY name", search)
	} else {
		rows, err = db.Query("SELECT id, name FROM Customer ORDER BY name")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*Customer, 0)
	for rows.Next() {
		result := new(Customer)
		err := rows.Scan(&result.Id, &result.Name)
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

func GetCustomer(db *sql.DB, id int) (*Customer, error) {
	row := db.QueryRow("SELECT id, name FROM Customer WHERE id = $1", id)
	result := new(Customer)
	err := row.Scan(&result.Id, &result.Name)
	if err != nil {
		return &Customer{}, err
	}
	return result, nil
}

func CreateCustomer(db *sql.DB, data Customer) (*Customer, error) {
	sqlStatement := `INSERT INTO Customer(name) VALUES ($1) RETURNING id`
	var id int
	err := db.QueryRow(sqlStatement, data.Name).Scan(&id)
	if err != nil {
		return &Customer{}, err
	}
	return GetCustomer(db, id)
}

func UpdateCustomer(db *sql.DB, id int, data Customer) (*Customer, error) {
	_, err := db.Exec("UPDATE Customer SET name = $1 WHERE id = $2", data.Name, id)
	if err != nil {
		return &Customer{}, err
	}
	return GetCustomer(db, id)
}

func DeleteCustomer(db *sql.DB, id int) error {
	if _, err := GetCustomer(db, id); err != nil {
		return err
	}

	result, err := db.Exec("DELETE FROM Customer WHERE id = $1", id)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}
