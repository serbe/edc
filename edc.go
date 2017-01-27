package edc

import (
	"log"
	"os"

	"gopkg.in/pg.v5"
)

// Edb struct to store *DB
type Edb struct {
	db  *pg.DB
	log bool
}

// SelectItem - struct for select element
type SelectItem struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// InitDB initialize database
func InitDB(dbname string, user string, password string, logsql bool) (*Edb, error) {
	e := new(Edb)
	opt := pg.Options{
		User:     user,
		Password: password,
		Database: dbname,
	}
	if logsql == true {
		pg.SetQueryLogger(log.New(os.Stdout, "", log.LstdFlags))
	}
	e.db = pg.Connect(&opt)
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

	return nil
}
