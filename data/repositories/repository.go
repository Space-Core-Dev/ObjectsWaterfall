package repositories

import (
	"fmt"

	"object-shooter.com/data"
)

const SQ_LITE = "sqlite3"

type Repository[T any] interface {
	SetData(tableName string, data T) error
	SetChankData(tableName string, jData []T) error
	GetData(tableName string, isRandom bool, take int, skip int64) ([]T, error)
	Count(tableName string) (int64, error)
}

func NewRepository[T any]() (Repository[T], error) {
	switch data.DbContext.Driver {
	case SQ_LITE:
		return MySqlRepositiry[T]{}, nil
	default:
		return nil, fmt.Errorf("there is no repository for %s driver", data.DbContext.Driver)
	}
}
