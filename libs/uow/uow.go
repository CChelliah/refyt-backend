package uow

import (
	"context"
	"database/sql"
	"fmt"
)

type (
	UnitOfWorkManager interface {
		Execute(ctx context.Context, uowHandler UnitOfWorkHandler) (err error)
	}

	UnitOfWork interface {
		GetTx() *sql.Tx
	}

	UnitOfWorkManagerImpl struct {
		db *sql.DB
	}

	UnitOfWorkImpl struct {
		tx *sql.Tx
	}

	UnitOfWorkHandler func(ctx context.Context, uow UnitOfWork) (err error)
)

func NewUnitOfWorkManager(db *sql.DB) *UnitOfWorkManagerImpl {
	return &UnitOfWorkManagerImpl{
		db: db,
	}
}

func (uow *UnitOfWorkManagerImpl) Execute(ctx context.Context, uowHandler UnitOfWorkHandler) (err error) {

	tx, err := uow.db.BeginTx(ctx, nil)

	if err != nil {
		return fmt.Errorf("unable to start uow handler transaction")
	}

	if err = uowHandler(ctx, &UnitOfWorkImpl{tx: tx}); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("unable to complete uow handler transaction, rolling back, %s", err.Error())
	}

	err = tx.Commit()

	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("unable to complete uow handler transaction, rolling back, %s", err.Error())
	}
	return nil
}

func (uow UnitOfWorkImpl) GetTx() *sql.Tx {
	return uow.tx
}
