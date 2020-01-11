package edc

import (
	"context"
	"time"
)

// Tcc - struct for tcc
type Tcc struct {
	ID        int64  `sql:"id"         json:"id"         form:"id"         query:"id"`
	Address   string `sql:"address"    json:"address"    form:"address"    query:"address"`
	ContactID int64  `sql:"contact_id" json:"contact_id" form:"contact_id" query:"contact_id"`
	CompanyID int64  `sql:"company_id" json:"company_id" form:"company_id" query:"company_id"`
	Note      string `sql:"note"       json:"note"       form:"note"       query:"note"`
	CreatedAt string `sql:"created_at" json:"-"`
	UpdatedAt string `sql:"updated_at" json:"-"`
}

// TccList - struct for tcc list
type TccList struct {
	ID        int64  `sql:"id"         json:"id"         form:"id"         query:"id"`
	Address   string `sql:"address"    json:"address"    form:"address"    query:"address"`
	ContactID int64  `sql:"contact_id" json:"contact_id" form:"contact_id" query:"contact_id"`
	Note      string `sql:"note"       json:"note"       form:"note"       query:"note"`
}

// TccGet - get one tcc by id
func TccGet(id int64) (Tcc, error) {
	var tcc Tcc
	if id == 0 {
		return tcc, nil
	}
	tcc.ID = id
	err := pool.QueryRow(context.Background(), `
		SELECT
			address,
			contact_id,
			company_id,
			note,
			created_at,
			updated_at
		FROM
			tccs
		WHERE
			id = $1
	`, id).Scan(&tcc.Address, &tcc.Address, &tcc.ContactID, &tcc.CompanyID, &tcc.Note, &tcc.CreatedAt, &tcc.UpdatedAt)
	if err != nil {
		errmsg("GetTcc select", err)
	}
	return tcc, err
}

// TccListGet - get all tcc for list
func TccListGet() ([]TccList, error) {
	var tccs []TccList
	rows, err := pool.Query(context.Background(), `
		SELECT
			id,
			address,
			contact_id,
			note,
		FROM
			tccs
		ORDER BY
			name ASC
	`)
	if err != nil {
		errmsg("TccListGet Query", err)
	}
	for rows.Next() {
		var tcc TccList
		err := rows.Scan(&tcc.ID, &tcc.Address, &tcc.ContactID, &tcc.Note)
		if err != nil {
			errmsg("TccListGet Scan", err)
			return tccs, err
		}
		tccs = append(tccs, tcc)
	}
	return tccs, rows.Err()
}

// TccInsert - create new tcc
func TccInsert(tcc Tcc) (int64, error) {
	err := pool.QueryRow(context.Background(), `
		INSERT INTO tccs
		(
			address,
			contact_id,
			company_id,
			note,
			created_at,
			updated_at
		)
		VALUES
		(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		)
		RETURNING
			id
	`, tcc.Address, tcc.ContactID, tcc.CompanyID, tcc.Note, time.Now(), time.Now()).Scan(&tcc.ID)
	if err != nil {
		errmsg("CreateTcc insert", err)
	}
	return tcc.ID, err
}

// TccUpdate - save tcc changes
func TccUpdate(tcc Tcc) error {
	_, err := pool.Exec(context.Background(), `
		UPDATE tccs SET
			address = $2,
			contact_id = $3,
			company_id = $4,
			note = $5,
			updated_at = $6
		WHERE
			id = $1
	`, tcc.ID, tcc.Address, tcc.ContactID, tcc.CompanyID, tcc.Note, time.Now())
	if err != nil {
		errmsg("UpdateTcc update", err)
	}
	return err
}

// TccDelete - delete tcc by id
func TccDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.Exec(context.Background(), `
		DELETE FROM
			tccs
		WHERE
			id = $1
	`, id)
	if err != nil {
		errmsg("DeleteTcc Exec", err)
	}
	return err
}

func tccCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			tccs (
				id         bigserial PRIMARY KEY,
				address    text,
				contact_id bigint,
				company_id bigint,
				note       text,
				created_at TIMESTAMP without time zone,
				updated_at TIMESTAMP without time zone default now(),
				UNIQUE(num_id, num_pass, type_id)
			)
	`
	_, err := pool.Exec(context.Background(), str)
	if err != nil {
		errmsg("tccCreateTable exec", err)
	}
	return err
}
