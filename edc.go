package edc

import (
	"log"
	"time"

	"github.com/go-pg/pg"
)

var logErrors bool

// Edb struct to store *DB
type Edb struct {
	db *pg.DB
}

// SelectItem - struct for select element
type SelectItem struct {
	ID   int64  `json:"id"   form:"id"   query:"id"`
	Name string `json:"name" form:"name" query:"name"`
}

// InitDB initialize database
func InitDB(
	host,
	port,
	dbname,
	user,
	password string,
	logsql,
	logerr bool,
) (*Edb, error) {
	e := new(Edb)
	opt := pg.Options{
		User:     user,
		Password: password,
		Database: dbname,
	}
	if host != "" || port != "" {
		opt.Addr = host + ":" + port
	}
	e.db = pg.Connect(&opt)
	if logsql {
		e.db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
			query, err := event.FormattedQuery()
			if err != nil {
				panic(err)
			}

			log.Printf("%s %s", time.Since(event.StartTime), query)
		})
	}
	logErrors = logerr
	err := e.createAllTables()

	return e, err
}

func (e *Edb) createAllTables() error {
	err := e.educationCreateTable()
	if err != nil {
		return err
	}
	err = e.kindCreateTable()
	if err != nil {
		return err
	}
	err = e.emailCreateTable()
	if err != nil {
		return err
	}
	err = e.companyCreateTable()
	if err != nil {
		return err
	}
	err = e.contactCreateTable()
	if err != nil {
		return err
	}
	err = e.postCreateTable()
	if err != nil {
		return err
	}
	err = e.rankCreateTable()
	if err != nil {
		return err
	}
	err = e.scopeCreateTable()
	if err != nil {
		return err
	}
	err = e.phoneCreateTable()
	if err != nil {
		return err
	}
	err = e.practiceCreateTable()
	if err != nil {
		return err
	}
	err = e.departmentCreateTable()
	if err != nil {
		return err
	}
	err = e.sirenTypeCreateTable()
	if err != nil {
		return err
	}
	err = e.sirenCreateTable()
	if err != nil {
		return err
	}
	err = e.certificateCreateTable()
	// if err != nil {
	// 	return err
	// }
	// err = e.tccCreateTable()
	return err
}
