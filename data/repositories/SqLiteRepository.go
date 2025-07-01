package repositories

import (
	"fmt"

	"objectswaterfall.com/data"
)

type mySqlRepositiry[T any] struct {
}

func (r mySqlRepositiry[T]) SetData(tableName string, jData T) error {
	if err := createTable(tableName); err != nil {
		return nil
	}

	stmt, err := data.DbContext.Db.Prepare(fmt.Sprintf(data.InsertData, tableName))
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(jData)
	if err != nil {
		return err
	}

	return nil
}

func (r mySqlRepositiry[T]) SetChankData(tableName string, jData []T) error {
	if err := createTable(tableName); err != nil {
		return err
	}
	tx, err := data.DbContext.Db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(fmt.Sprintf(data.InsertData, tableName))
	if err != nil {
		return err
	}

	defer stmt.Close()
	for _, v := range jData {
		if _, err := stmt.Exec(v); err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r mySqlRepositiry[T]) GetData(tableName string, isRandom bool, take int, skip int64) ([]T, error) {
	rows, err := data.DbContext.Db.Query(fmt.Sprintf(data.GetJson, tableName), skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jsons []T
	for rows.Next() {
		var data T
		err := rows.Scan(&data)
		if err != nil {
			return nil, err
		}
		jsons = append(jsons, data)
	}

	return jsons, nil
}

func (r mySqlRepositiry[T]) Count(tableName string) (int64, error) {
	var count int64
	err := data.DbContext.Db.QueryRow(fmt.Sprintf(data.Count, tableName)).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r mySqlRepositiry[T]) GetAllTables() ([]string, error) {

	rows, err := data.DbContext.Db.Query(data.Tables)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		names = append(names, name)
	}

	return names, nil
}

func createTable(tableName string) error {
	stmt, err := data.DbContext.Db.Prepare(fmt.Sprintf(data.CreateTable, tableName))
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}
