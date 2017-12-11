package category

import (
	"database/sql"
)

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func ListCategories(db *sql.DB, search string) ([]*Category, error) {
	var rows *sql.Rows
	var err error

	if search != "" {
		rows, err = db.Query("SELECT id, name FROM Category WHERE UPPER(name) LIKE CONCAT('%', UPPER($1), '%') ORDER BY name", search)
	} else {
		rows, err = db.Query("SELECT id, name FROM Category ORDER BY name")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*Category, 0)
	for rows.Next() {
		result := new(Category)
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

func GetCategory(db *sql.DB, id int) (*Category, error) {
	row := db.QueryRow("SELECT id, name FROM Category WHERE id = $1", id)
	result := new(Category)
	err := row.Scan(&result.Id, &result.Name)
	if err != nil {
		return &Category{}, err
	}
	return result, nil
}

func CreateCategory(db *sql.DB, data Category) (*Category, error) {
	sqlStatement := `INSERT INTO Category(name) VALUES ($1) RETURNING id`
	var id int
	err := db.QueryRow(sqlStatement, data.Name).Scan(&id)
	if err != nil {
		return &Category{}, err
	}
	return GetCategory(db, id)
}

func UpdateCategory(db *sql.DB, id int, data Category) (*Category, error) {
	_, err := db.Exec("UPDATE Category SET name = $1 WHERE id = $2", data.Name, id)
	if err != nil {
		return &Category{}, err
	}
	return GetCategory(db, id)
}

func DeleteCategory(db *sql.DB, id int) error {
	if _, err := GetCategory(db, id); err != nil {
		return err
	}

	result, err := db.Exec("DELETE FROM Category WHERE id = $1", id)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}
