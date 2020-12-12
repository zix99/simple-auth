package appcontext

import (
	"simple-auth/pkg/db"

	"github.com/labstack/echo/v4"
)

const dbContextKey = "appcontext.sadb"

type sadbWrapper struct {
	root        db.SADB
	transaction db.SADBTransaction
}

func WithSADB(db db.SADB) ProviderFunc {
	return func(c Context) (string, interface{}) {
		log := GetLogger(c)
		return dbContextKey, &sadbWrapper{
			root: db.WithLogger(log),
		}
	}
}

func GetSADB(c Context) db.SADB {
	dbw := c.Get(dbContextKey).(*sadbWrapper)
	if dbw.transaction != nil {
		return dbw.transaction
	}
	return dbw.root
}

func Transaction() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			dbw := c.Get(dbContextKey).(*sadbWrapper)
			log := GetLogger(c)

			if dbw.transaction != nil {
				panic("Nested transactions unsupported")
			}

			dbw.transaction = dbw.root.BeginTransaction()

			committed := false
			defer func() {
				if !committed {
					log.Warnln("Rollingback transaction")
					if dberr := dbw.transaction.Rollback(); dberr != nil {
						log.Errorln(dberr)
					}
				}
			}()

			err := next(c)

			if err == nil {
				log.Debugln("Committing transaction...")
				if dberr := dbw.transaction.Commit(); dberr != nil {
					log.Errorln(dberr)
				} else {
					committed = true
				}
			}

			return err
		}
	}
}
