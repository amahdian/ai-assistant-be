package storage

import (
	"gorm.io/gorm/schema"
)

// CrudStorage is the base storage class that provides common functionalities which all stores can benefit from.
type CrudStorage[M schema.Tabler] interface {
	CreateOne(model M) error
	CreateMany(models []M) error
	CreateManyWithAssociation(models []M, saveAssociations bool) error
	CreateInBatches(models []M) error

	FindById(id string) (model M, err error)
	ListByIds(ids []string) (models []M, err error)

	UpdateOne(model M, updateZeroValues bool) error
	UpdateMany(models []M) error

	UpsertOne(model M, saveAssociations bool) error
	UpdatePartial(model M, returnUpdated bool) error
	UpsertMany(models []M) error
	UpsertManyWithAssociation(models []M, saveAssociations bool) error
	UpsertInBatches(models []M) error

	ExistsById(id string) (exists bool, err error)

	DeleteOne(model M) error
	DeleteMany(models []M) error
	DeleteById(id string) error
	DeleteByIds(ids []string) error

	ListAll() (models []M, err error)
}
