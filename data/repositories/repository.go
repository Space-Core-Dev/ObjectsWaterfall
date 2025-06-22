package repositories

import (
	"fmt"

	"object-shooter.com/data"
)

const SQ_LITE = "sqlite3"

type Repository interface {
	SetData(tableName string, data string) error
	SetChankData(tableName string, jData []string) error
	GetData(tableName string, take, skip int) string
}

func NewRepository() (Repository, error) {
	switch data.DbContext.Driver {
	case SQ_LITE:
		return &MySqlRepositiry{}, nil
	default:
		return nil, fmt.Errorf("there is no repository for %s driver", data.DbContext.Driver)
	}
}
