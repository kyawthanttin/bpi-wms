package dbutil

import (
	"bytes"
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const MaxResults = 500
const MaxOpenConns = 100
const MaxIdleConns = 20

// Return the postgresql connection
func NewDB(datasourceName string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", datasourceName)
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetMaxOpenConns(MaxOpenConns)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Check for the provided value in the database records. Return true if found
func IsExist(db *sqlx.DB, tbl string, col string, val interface{}) (bool, error) {
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM "+tbl+" WHERE "+col+" = $1", val)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

// Check for the provided value in the database records except the record specified by id. Return true if found
func IsExistExcept(db *sqlx.DB, tbl string, id int, col string, val interface{}) (bool, error) {
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM "+tbl+" WHERE id != $1 AND "+col+" = $2", id, val)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func Insert(db *sqlx.DB, tbl string, entity interface{}) (interface{}, error) {
	cols, vals, err := extractInsertFields(entity)
	if err != nil {
		return nil, err
	}
	colBuff := bytes.NewBufferString("")
	valBuff := bytes.NewBufferString("")
	colLen := len(cols)
	for i := 1; i <= colLen; i++ {
		colBuff.WriteString(cols[i-1])
		valBuff.WriteString("$" + strconv.Itoa(i))
		if i < colLen {
			colBuff.WriteString(", ")
			valBuff.WriteString(", ")
		}
	}
	sqlStmt := "INSERT INTO " + tbl + "(" + colBuff.String() + ") VALUES (" + valBuff.String() + ") RETURNING id"
	var id int
	err = db.QueryRow(sqlStmt, vals...).Scan(&id)
	if err != nil {
		return nil, err
	}
	return id, err
}

func Update(db *sqlx.DB, tbl string, entity interface{}, query interface{}) error {
	cols, vals, err := extractUpdateFields(entity)
	if err != nil {
		return err
	}
	qcols, qvals, err := extractQueryFields(query)
	if err != nil {
		return err
	}
	colBuff := bytes.NewBufferString("")
	colLen := len(cols)
	for i := 1; i <= colLen; i++ {
		colBuff.WriteString(cols[i-1] + " = $" + strconv.Itoa(i))
		if i < colLen {
			colBuff.WriteString(", ")
		}
	}
	qBuff := bytes.NewBufferString("")
	qcolLen := len(qcols)
	for i := 1; i <= qcolLen; i++ {
		qBuff.WriteString(qcols[i-1] + " = $" + strconv.Itoa(i+colLen))
		if i < qcolLen {
			qBuff.WriteString(" AND ")
		}
	}

	sqlStmt := "UPDATE " + tbl + " SET " + colBuff.String() + " WHERE " + qBuff.String()
	args := append(vals, qvals...)
	_, err = db.Exec(sqlStmt, args...)
	if err != nil {
		return err
	}
	return nil
}

func extractInsertFields(entity interface{}) ([]string, []interface{}, error) {
	return extractFields(entity, "i")
}

func extractUpdateFields(entity interface{}) ([]string, []interface{}, error) {
	return extractFields(entity, "u")
}

func extractQueryFields(entity interface{}) ([]string, []interface{}, error) {
	return extractFields(entity, "q")
}

func extractFields(entity interface{}, dbop string) ([]string, []interface{}, error) {
	ev := reflect.ValueOf(entity)
	if ev.Kind() != reflect.Ptr || ev.IsNil() {
		return nil, nil, errors.New("Pointer expected")
	}
	et := reflect.Indirect(ev).Type()
	cols := []string{}
	vals := make([]interface{}, 0)
	for i := 0; i < et.NumField(); i++ {
		v := et.Field(i)

		if dbop == "i" || dbop == "u" {
			// Check whether the field is bound for insert/update
			op := v.Tag.Get("dbop")
			if op == "" {
				continue
			} else if strings.Index(op, dbop) < 0 {
				continue
			}
		}

		el := reflect.Indirect(ev.Elem().FieldByName(v.Name))
		isString := false
		if dbop == "u" || dbop == "q" {
			// Skip if there is no value
			if !el.IsValid() {
				continue
			}

			switch el.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if el.Int() == 0 {
					continue
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if el.Uint() == 0 {
					continue
				}
			case reflect.Float32, reflect.Float64:
				if el.Float() == 0 {
					continue
				}
			case reflect.String:
				if el.String() == "" {
					continue
				}
				isString = true

			default:
				continue
			}
		}

		colname := strings.ToLower(v.Name)
		dbcol := v.Tag.Get("db")
		if dbcol != "" {
			cols = append(cols, dbcol)
		} else {
			cols = append(cols, colname)
		}
		if isString {
			vals = append(vals, strings.TrimSpace(el.String()))
		} else {
			vals = append(vals, el.Interface())
		}
	}
	if len(cols) <= 0 {
		return nil, nil, errors.New("No eligible fields")
	}
	return cols, vals, nil
}
