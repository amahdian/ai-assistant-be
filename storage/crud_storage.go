package storage

import (
	"gorm.io/gorm/schema"
)

// CrudStorage is the base storage class that provides common functionalities which all stores can benefit from.
// Please add your common storage logic here.
type CrudStorage[M schema.Tabler] interface {
	CreateOne(model M) error
	CreateMany(models []M) error

	FindById(id int64) (model M, err error)
	ListByIds(ids []int64) (models []M, err error)

	UpdateOne(model M, saveAssociations bool) error
	UpdatePartial(model M, returnUpdated bool) error
	UpdateMany(models []M) error

	ExistsById(id int64) (exists bool, err error)
	DeleteById(id int64) error
	DeleteByIds(ids []int64) error

	ListAll() (models []M, err error)
}
