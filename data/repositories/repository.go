package repositories

import (
	"fmt"

	"objectswaterfall.com/data"
)

const SQ_LITE = "sqlite3"

type Repository[T any] interface {
	SetData(tableName string, data T) error
	SetChankData(tableName string, jData []T) error
	GetData(tableName string, isRandom bool, take int, skip int64) ([]T, error)
	Count(tableName string) (int64, error)
}

type SqLiteRepository[T any] interface {
	Repository[T]
	GetAllTables() ([]string, error)
}

func NewRepository[T any]() (SqLiteRepository[T], error) {
	switch data.DbContext.Driver {
	case SQ_LITE:
		return mySqlRepositiry[T]{}, nil
	default:
		return nil, fmt.Errorf("there is no repository for %s driver", data.DbContext.Driver)
	}
}
