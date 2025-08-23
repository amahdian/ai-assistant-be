package pg

import (
	"context"
	"fmt"
	"reflect"

	"github.com/amahdian/ai-assistant-be/storage"
	"gorm.io/gorm"
)

type Stg struct {
	db *gorm.DB
}

func NewStg(db *gorm.DB) storage.Storage {
	return &Stg{db: db}
}

func (stg *Stg) WithContext(ctx context.Context) storage.Storage {
	return &Stg{
		db: stg.mustOrmSession(ctx).db,
	}
}

func (stg *Stg) Atomic(fn func(atomicStorage storage.Storage) error) (err error) {
	tx := stg.db.Begin()
	if tx.Error != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
		if err != nil {
			if rbErr := tx.Rollback().Error; rbErr != nil {
				err = fmt.Errorf("Transaction err: %v. something went wrong in rollback: %v", err, rbErr)
			}
		} else {
			err = tx.Commit().Error
		}
	}()

	txStorage := &Stg{
		db: tx,
	}

	err = fn(txStorage)
	return
}

func (stg *Stg) RegisterDeleteHook(fn func(ctx context.Context, storage storage.Storage, entity interface{}) error) {
	_ = stg.db.Callback().Delete().Before("gorm:delete").Register("app_delete_hook", func(db *gorm.DB) {
		if db.Error != nil {
			return
		}

		switch db.Statement.ReflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			db.Statement.CurDestIndex = 0
			for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
				if value := reflect.Indirect(db.Statement.ReflectValue.Index(i)); value.CanAddr() {
					_ = db.AddError(fn(db.Statement.Context, stg, value.Addr().Interface()))
				} else {
					_ = db.AddError(gorm.ErrInvalidValue)
					return
				}
				db.Statement.CurDestIndex++
			}
		case reflect.Struct:
			if db.Statement.ReflectValue.CanAddr() {
				_ = db.AddError(fn(db.Statement.Context, stg, db.Statement.ReflectValue.Addr().Interface()))
			} else {
				_ = db.AddError(gorm.ErrInvalidValue)
			}
		}
	})
}

func (stg *Stg) Session(ctx context.Context) storage.Session {
	return stg.mustOrmSession(ctx)
}

func (stg *Stg) Begin(ctx context.Context) (context.Context, storage.Session, error) {
	session, err := stg.Session(ctx).Begin()
	if err != nil {
		return nil, nil, err
	}
	ctx = storage.WithContext(ctx, session)
	return ctx, session, nil
}

func (stg *Stg) mustOrmSession(ctx context.Context) *ormSession {
	if ses := storage.FromContext(ctx); ses != nil {
		return ses.(*ormSession)
	}
	db := stg.db.WithContext(ctx)
	ses := &ormSession{db: db}
	return ses
}

func (stg *Stg) User(ctx context.Context) storage.UserStorage {
	return NewUserStg(stg.mustOrmSession(ctx))
}
