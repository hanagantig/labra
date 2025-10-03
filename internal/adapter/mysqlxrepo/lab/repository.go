package lab

import (
	"github.com/jmoiron/sqlx"
	"labra/internal/adapter/mysqlxrepo"
)

type Repository struct {
	mysqlxrepo.Transactor
}

func NewRepository(conn *sqlx.DB) *Repository {
	return &Repository{
		mysqlxrepo.NewTransactor(conn),
	}
}
