package edc

import (
	"context"
	"time"
)

// Kind - struct for kind
type Kind struct {
	ID        int64  `sql:"id"         json:"id"         form:"id"         query:"id"`
	Name      string `sql:"name"       json:"name"       form:"name"       query:"name"`
	ShortName string `sql:"short_name" json:"short_name" form:"short_name" query:"short_name"`
	Note      string `sql:"note"       json:"note"       form:"note"       query:"note"`
	CreatedAt string `sql:"created_at" json:"-"`
	UpdatedAt string `sql:"updated_at" json:"-"`
}

// KindList - struct for kind list
type KindList struct {
	ID        int64  `sql:"id"         json:"id"         form:"id"         query:"id"`
	Name      string `sql:"name"       json:"name"       form:"name"       query:"name"`
	ShortName string `sql:"short_name" json:"short_name" form:"short_name" query:"short_name"`
	Note      string `sql:"note"       json:"note"       form:"note"       query:"note"`
}

// KindGet - get one kind by id
func KindGet(id int64) (Kind, error) {
	var kind Kind
	if id == 0 {
		return kind, nil
	}
	kind.ID = id
	err := pool.QueryRow(context.Background(), `
		SELECT
			name,
			short_name,
			note,
			created_at,
			updated_at
		FROM
			kinds
		WHERE
			id = $1
	`, id).Scan(&kind.Name, &kind.ShortName, &kind.Note, &kind.CreatedAt, &kind.UpdatedAt)
	if err != nil {
		errmsg("KindGet QueryRow", err)
	}
	return kind, err
}

// KindListGet - get all kind for list
func KindListGet() ([]KindList, error) {
	var kinds []KindList
	rows, err := pool.Query(context.Background(), `
		SELECT
			id,
			name,
			short_name,
			note
		FROM
			kinds
		ORDER BY
			name ASC
	`)
	if err != nil {
		errmsg("KindListGet Query", err)
	}
	for rows.Next() {
		var kind KindList
		err := rows.Scan(&kind.ID, &kind.Name, &kind.ShortName, &kind.Note)
		if err != nil {
			errmsg("KindListGet Scan", err)
			return kinds, err
		}
		kinds = append(kinds, kind)
	}
	return kinds, rows.Err()
}

// KindSelectGet - get all kind for select
func KindSelectGet() ([]SelectItem, error) {
	var kinds []SelectItem
	rows, err := pool.Query(context.Background(), `
		SELECT
			id,
			name
		FROM
			kinds
		ORDER BY
			name ASC
	`)
	if err != nil {
		errmsg("KindSelectGet Query", err)
	}
	for rows.Next() {
		var kind SelectItem
		err := rows.Scan(&kind.ID, &kind.Name)
		if err != nil {
			errmsg("KindSelectGet Scan", err)
			return kinds, err
		}
		kinds = append(kinds, kind)
	}
	return kinds, rows.Err()
}

// KindInsert - create new kind
func KindInsert(kind Kind) (int64, error) {
	err := pool.QueryRow(context.Background(), `
		INSERT INTO educations
		(
			name,
			short_name,
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
	`, kind.Name, kind.ShortName, kind.Note, time.Now(), time.Now()).Scan(&kind.ID)
	if err != nil {
		errmsg("KindInsert QueryRow", err)
	}
	return kind.ID, nil
}

// KindUpdate - save kind changes
func KindUpdate(kind Kind) error {
	_, err := pool.Exec(context.Background(), `
		UPDATE educations SET
			name = $2,
			short_name = $3,
			note = $4,
			updated_at = $5
		WHERE
			id = $1
	`, kind.ID, kind.Name, kind.ShortName, kind.Note, time.Now())
	if err != nil {
		errmsg("KindUpdate Exec", err)
	}
	return err
}

// KindDelete - delete kind by id
func KindDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.Exec(context.Background(), `
		DELETE FROM
			kinds
		WHERE
			id = $1
	`, id)
	if err != nil {
		errmsg("DeleteKind Exec", err)
	}
	return err
}

func kindCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			kinds (
				id bigserial primary key,
				name text,
				short_name text,
				note text,
				created_at TIMESTAMP without time zone,
				updated_at TIMESTAMP without time zone default now(),
				UNIQUE(name)
			)
	`
	_, err := pool.Exec(context.Background(), str)
	if err != nil {
		errmsg("kindCreateTable exec", err)
	}
	return err
}
