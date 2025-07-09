package repositories

import (
	"database/sql"
	"fmt"

	"objectswaterfall.com/core/models"
	"objectswaterfall.com/data"
)

type mySqlRepositiry[T any] struct {
}

func (r mySqlRepositiry[T]) SetData(workerName string, jData T) error {
	if err := createTable(workerName); err != nil {
		return nil
	}

	stmt, err := data.DbContext.Db.Prepare(fmt.Sprintf(data.InsertData, workerName))
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

func (r mySqlRepositiry[T]) SetChankData(workerName string, jData []T) error {
	if err := createTable(workerName); err != nil {
		return err
	}
	tx, err := data.DbContext.Db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(fmt.Sprintf(data.InsertData, workerName))
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

func (r mySqlRepositiry[T]) GetData(workerName string, isRandom bool, take int, skip int64) ([]T, error) {
	rows, err := data.DbContext.Db.Query(fmt.Sprintf(data.GetJson, workerName), skip, take)
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

func (r mySqlRepositiry[T]) Count(workerName string) (int64, error) {
	var count int64
	err := data.DbContext.Db.QueryRow(fmt.Sprintf(data.Count, workerName)).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r mySqlRepositiry[T]) GetAllWorkers() ([]string, error) {
	rows, err := data.DbContext.Db.Query(data.Workers)
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

func (r mySqlRepositiry[T]) AddSettings(settings models.BackgroundWorkerSettings) error {
	var err error
	var existingTable string
	err = data.DbContext.Db.QueryRow(data.Exists, settings.WorkerName).Scan(&existingTable)
	if err != nil && err != sql.ErrNoRows {
		return err
	} else if existingTable != "" {
		return fmt.Errorf("table %s already exists", settings.WorkerName)
	}

	_, err = data.DbContext.Db.Exec(data.CreateWorkerSettingsTable)
	if err != nil {
		return err
	}

	stmt, err := data.DbContext.Db.Prepare(data.InsertWorkerSettings)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		settings.WorkerName,
		settings.Timer,
		settings.RequestDelay,
		settings.Random,
		settings.WritesNumberToSend,
		settings.TotalToSend,
		settings.StopWhenTableEnds)
	if err != nil {
		return err
	}

	err = createTable(settings.WorkerName)
	if err != nil {
		return err
	}

	return nil
}

func (r mySqlRepositiry[T]) GetWorkerSettings(settingsWorkerName string) (*models.BackgroundWorkerSettings, error) {
	row := data.DbContext.Db.QueryRow(data.GetWorkerSettings, settingsWorkerName)

	var settings models.BackgroundWorkerSettings
	err := row.Scan(&settings.WorkerName,
		&settings.Timer,
		&settings.RequestDelay,
		&settings.Random,
		&settings.WritesNumberToSend,
		&settings.TotalToSend,
		&settings.StopWhenTableEnds)

	if err != nil {
		return nil, err
	}

	return &settings, nil
}

func createTable(workerName string) error {
	stmt, err := data.DbContext.Db.Prepare(fmt.Sprintf(data.CreateTable, workerName))
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
