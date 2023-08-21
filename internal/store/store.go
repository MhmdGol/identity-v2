package store

import (
	"database/sql"
	"fmt"
	"identity-v2/cmd/config"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mssqldialect"
)

func NewSQLStorage(conf config.Config) (*bun.DB, error) {
	connString := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.Name,
	)

	sqldb, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, err
	}

	db := bun.NewDB(sqldb, mssqldialect.New())

	return db, nil
}
