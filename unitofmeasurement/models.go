package unitofmeasurement

import (
	"database/sql"
)

type UnitOfMeasurement struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Abbr string `json:"abbr"`
}

func ListUnitOfMeasurements(db *sql.DB, search string) ([]*UnitOfMeasurement, error) {
	var rows *sql.Rows
	var err error

	if search != "" {
		rows, err = db.Query("SELECT id, name, abbr FROM Unit_Of_Measurement WHERE UPPER(name) LIKE CONCAT('%', UPPER($1), '%') ORDER BY name", search)
	} else {
		rows, err = db.Query("SELECT id, name, abbr FROM Unit_Of_Measurement ORDER BY name")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*UnitOfMeasurement, 0)
	for rows.Next() {
		result := new(UnitOfMeasurement)
		err := rows.Scan(&result.Id, &result.Name, &result.Abbr)
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

func GetUnitOfMeasurement(db *sql.DB, id int) (*UnitOfMeasurement, error) {
	row := db.QueryRow("SELECT id, name, abbr FROM Unit_Of_Measurement WHERE id = $1", id)
	result := new(UnitOfMeasurement)
	err := row.Scan(&result.Id, &result.Name, &result.Abbr)
	if err != nil {
		return &UnitOfMeasurement{}, err
	}
	return result, nil
}

func CreateUnitOfMeasurement(db *sql.DB, data UnitOfMeasurement) (*UnitOfMeasurement, error) {
	sqlStatement := `INSERT INTO Unit_Of_Measurement(name, abbr) VALUES ($1) RETURNING id`
	var id int
	err := db.QueryRow(sqlStatement, data.Name, data.Abbr).Scan(&id)
	if err != nil {
		return &UnitOfMeasurement{}, err
	}
	return GetUnitOfMeasurement(db, id)
}

func UpdateUnitOfMeasurement(db *sql.DB, id int, data UnitOfMeasurement) (*UnitOfMeasurement, error) {
	_, err := db.Exec("UPDATE Unit_Of_Measurement SET name = $1, abbr = $2 WHERE id = $3", data.Name, data.Abbr, id)
	if err != nil {
		return &UnitOfMeasurement{}, err
	}
	return GetUnitOfMeasurement(db, id)
}

func DeleteUnitOfMeasurement(db *sql.DB, id int) error {
	if _, err := GetUnitOfMeasurement(db, id); err != nil {
		return err
	}

	result, err := db.Exec("DELETE FROM Unit_Of_Measurement WHERE id = $1", id)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}
