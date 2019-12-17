package edc

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	logErrors bool
	pool      *pgxpool.Pool
)

// SelectItem - struct for select element
type SelectItem struct {
	ID   int64  `json:"id"   form:"id"   query:"id"`
	Name string `json:"name" form:"name" query:"name"`
}

// InitDB initialize database
func InitDB(
	db_url string,
	logsql,
	logerr bool,
) error {
	var err error
	pool, err = pgxpool.Connect(context.Background(), db_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}
	logErrors = logerr
	return AllTablesInsert()
}

func AllTablesInsert() error {
	err := educationCreateTable()
	if err != nil {
		return err
	}
	err = kindCreateTable()
	if err != nil {
		return err
	}
	err = emailCreateTable()
	if err != nil {
		return err
	}
	err = companyCreateTable()
	if err != nil {
		return err
	}
	err = contactCreateTable()
	if err != nil {
		return err
	}
	err = postCreateTable()
	if err != nil {
		return err
	}
	err = rankCreateTable()
	if err != nil {
		return err
	}
	err = scopeCreateTable()
	if err != nil {
		return err
	}
	err = phoneCreateTable()
	if err != nil {
		return err
	}
	err = practiceCreateTable()
	if err != nil {
		return err
	}
	err = departmentCreateTable()
	if err != nil {
		return err
	}
	err = sirenTypeCreateTable()
	if err != nil {
		return err
	}
	err = sirenCreateTable()
	if err != nil {
		return err
	}
	err = certificateCreateTable()
	// if err != nil {
	// 	return err
	// }
	// err = e.tccCreateTable()
	return err
}
