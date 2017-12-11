package country

import (
	"database/sql"
)

type Country struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func ListCountries(db *sql.DB, search string) ([]*Country, error) {
	var rows *sql.Rows
	var err error

	if search != "" {
		rows, err = db.Query("SELECT id, name FROM Country WHERE UPPER(name) LIKE CONCAT('%', UPPER($1), '%') ORDER BY name", search)
	} else {
		rows, err = db.Query("SELECT id, name FROM Country ORDER BY name")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*Country, 0)
	for rows.Next() {
		result := new(Country)
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

func GetCountry(db *sql.DB, id int) (*Country, error) {
	row := db.QueryRow("SELECT id, name FROM Country WHERE id = $1", id)
	result := new(Country)
	err := row.Scan(&result.Id, &result.Name)
	if err != nil {
		return &Country{}, err
	}
	return result, nil
}

func CreateCountry(db *sql.DB, data Country) (*Country, error) {
	sqlStatement := `INSERT INTO Country(name) VALUES ($1) RETURNING id`
	var id int
	err := db.QueryRow(sqlStatement, data.Name).Scan(&id)
	if err != nil {
		return &Country{}, err
	}
	return GetCountry(db, id)
}

func UpdateCountry(db *sql.DB, id int, data Country) (*Country, error) {
	_, err := db.Exec("UPDATE Country SET name = $1 WHERE id = $2", data.Name, id)
	if err != nil {
		return &Country{}, err
	}
	return GetCountry(db, id)
}

func DeleteCountry(db *sql.DB, id int) error {
	if _, err := GetCountry(db, id); err != nil {
		return err
	}

	result, err := db.Exec("DELETE FROM Country WHERE id = $1", id)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}
