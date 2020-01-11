package edc

import (
	"context"
	"time"
)

// SirenType - struct for sirenType
type SirenType struct {
	ID        int64  `sql:"id"         json:"id"            form:"id"     query:"id"`
	Name      string `sql:"name"       json:"name"          form:"name"   query:"name"`
	Radius    int64  `sql:"radius"     json:"radius,string" form:"radius" query:"radius"`
	Note      string `sql:"note"       json:"note"          form:"note"   query:"note"`
	CreatedAt string `sql:"created_at" json:"-"`
	UpdatedAt string `sql:"updated_at" json:"-"`
}

// SirenTypeList - struct for sirenType list
type SirenTypeList struct {
	ID     int64  `sql:"id"     json:"id"            form:"id"     query:"id"`
	Name   string `sql:"name"   json:"name"          form:"name"   query:"name"`
	Radius int64  `sql:"radius" json:"radius,string" form:"radius" query:"radius"`
	Note   string `sql:"note"   json:"note"          form:"note"   query:"note"`
}

// SirenTypeGet - get one sirenType by id
func SirenTypeGet(id int64) (SirenType, error) {
	var sirenType SirenType
	if id == 0 {
		return sirenType, nil
	}
	sirenType.ID = id
	err := pool.QueryRow(context.Background(), `
		SELECT
			name,
			radius,
			note,
			created_at,
			updated_at
		FROM
			siren_types
		WHERE
			id = $1
	`, id).Scan(&sirenType.Name, &sirenType.Radius, &sirenType.Note, &sirenType.CreatedAt, &sirenType.UpdatedAt)
	if err != nil {
		errmsg("SirenTypeGet QueryRow", err)
	}
	return sirenType, err
}

// SirenTypeListGet - get all sirenType for list
func SirenTypeListGet() ([]SirenTypeList, error) {
	var sirenTypes []SirenTypeList
	rows, err := pool.Query(context.Background(), `
		SELECT
			id,
			name,
			radius,
			note
		FROM
			siren_types
		ORDER BY
			name ASC
	`)
	if err != nil {
		errmsg("SirenTypeListGet Query", err)
	}
	for rows.Next() {
		var sirenType SirenTypeList
		err := rows.Scan(&sirenType.ID, &sirenType.Name, &sirenType.Radius, &sirenType.Note)
		if err != nil {
			errmsg("SirenTypeListGet Scan", err)
			return sirenTypes, err
		}
		sirenTypes = append(sirenTypes, sirenType)
	}
	return sirenTypes, rows.Err()
}

// SirenTypeSelectGet - get all sirenType for select
func SirenTypeSelectGet() ([]SelectItem, error) {
	var sirenTypes []SelectItem
	rows, err := pool.Query(context.Background(), `
		SELECT
			id,
			name
		FROM
			siren_types
		ORDER BY
			name ASC
	`)
	if err != nil {
		errmsg("SirenTypeSelectGet Query", err)
	}
	for rows.Next() {
		var sirenType SelectItem
		err := rows.Scan(&sirenType.ID, &sirenType.Name)
		if err != nil {
			errmsg("SirenTypeSelectGet Scan", err)
			return sirenTypes, err
		}
		sirenTypes = append(sirenTypes, sirenType)
	}
	return sirenTypes, rows.Err()
}

// SirenTypeInsert - create new sirenType
func SirenTypeInsert(sirenType SirenType) (int64, error) {
	err := pool.QueryRow(context.Background(), `
		INSERT INTO siren_types
		(
			name,
			radius,
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
			$5
		)
		RETURNING
			id
	`, sirenType.Name, sirenType.Radius, sirenType.Note, time.Now(), time.Now()).Scan(&sirenType.ID)
	if err != nil {
		errmsg("SirenTypeInsert QueryRow", err)
	}
	return sirenType.ID, nil
}

// SirenTypeUpdate - save sirenType changes
func SirenTypeUpdate(sirenType SirenType) error {
	_, err := pool.Exec(context.Background(), `
		UPDATE siren_types SET
			name = $2,
			radius = $3,
			note = $4,
			updated_at = $5
		WHERE
			id = $1
	`, sirenType.Name, sirenType.Radius, sirenType.Note, time.Now())
	if err != nil {
		errmsg("SirenTypeUpdate Exec", err)
	}
	return err
}

// SirenTypeDelete - delete sirenType by id
func SirenTypeDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.Exec(context.Background(), `
		DELETE FROM
			siren_types
		WHERE
			id = $1
	`, id)
	if err != nil {
		errmsg("SirenTypeDelete Exec", err)
	}
	return err
}

func sirenTypeCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			siren_types (
				id         bigserial primary key,
				name       text,
				radius     bigint,
				note       text,
				created_at TIMESTAMP without time zone,
				updated_at TIMESTAMP without time zone default now(),
				UNIQUE(name, radius)
			);`
	_, err := pool.Exec(context.Background(), str)
	if err != nil {
		errmsg("sirenCreateTable exec", err)
	}
	return err
}
