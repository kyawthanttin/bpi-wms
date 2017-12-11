package models

import (
	"time"
)

type Category struct {
	Id           int64     `json:"id"`
	Name         string    `json:"name"`
	Created      time.Time `json:"created"`
	Last_Updated time.Time `json:"lastUpdated"`
}

func (db *DB) AllCategories() ([]*Category, error) {
	rows, err := db.Query("SELECT * FROM Category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*Category, 0)
	for rows.Next() {
		result := new(Category)
		err := rows.Scan(&result.Id, &result.Name, &result.Created, &result.Last_Updated)
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
