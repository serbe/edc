package edc

import (
	"context"
	"time"
)

// Rank - struct for rank
type Rank struct {
	ID        int64  `sql:"id"         json:"id"   form:"id"   query:"id"`
	Name      string `sql:"name"       json:"name" form:"name" query:"name"`
	Note      string `sql:"note"       json:"note" form:"note" query:"note"`
	CreatedAt string `sql:"created_at" json:"-"`
	UpdatedAt string `sql:"updated_at" json:"-"`
}

// RankList - struct for rank list
type RankList struct {
	ID   int64  `sql:"id"   json:"id"   form:"id"   query:"id"`
	Name string `sql:"name" json:"name" form:"name" query:"name"`
	Note string `sql:"note" json:"note" form:"note" query:"note"`
}

// RankGet - get one rank by id
func RankGet(id int64) (Rank, error) {
	var rank Rank
	if id == 0 {
		return rank, nil
	}
	rank.ID = id
	err := pool.QueryRow(context.Background(), `
		SELECT
			name,
			note,
			created_at,
			updated_at
		FROM
			ranks
		WHERE
			id = $1
	`, id).Scan(&rank.Name, &rank.Note, &rank.CreatedAt, &rank.UpdatedAt)
	if err != nil {
		errmsg("RankGet QueryRow", err)
	}
	return rank, err
}

// RankListGet - get all rank for list
func RankListGet() ([]RankList, error) {
	var ranks []RankList
	rows, err := pool.Query(context.Background(), `
		SELECT
			id,
			name,
			note
		FROM
			ranks
		ORDER BY
			name ASC
	`)
	if err != nil {
		errmsg("RankListGet Query", err)
	}
	for rows.Next() {
		var rank RankList
		err := rows.Scan(&rank.ID, &rank.Name, &rank.Note)
		if err != nil {
			errmsg("PostListGet Scan", err)
			return ranks, err
		}
		ranks = append(ranks, rank)
	}
	return ranks, rows.Err()
}

// RankSelectGet - get all rank for select
func RankSelectGet() ([]SelectItem, error) {
	var ranks []SelectItem
	rows, err := pool.Query(context.Background(), `
		SELECT
			id,
			name
		FROM
			ranks
		ORDER BY
			name ASC
	`)
	if err != nil {
		errmsg("RankSelectGet Query", err)
	}
	for rows.Next() {
		var rank SelectItem
		err := rows.Scan(&rank.ID, &rank.Name)
		if err != nil {
			errmsg("RankSelectGet Scan", err)
			return ranks, err
		}
		ranks = append(ranks, rank)
	}
	return ranks, rows.Err()
}

// RankInsert - create new rank
func RankInsert(rank Rank) (int64, error) {
	err := pool.QueryRow(context.Background(), `
		INSERT INTO ranks
		(
			name,
			note,
			created_at,
			updated_at
		)
		VALUES
		(
			$1,
			$2,
			$3,
			$4
		)
		RETURNING
			id
	`, rank.Name, rank.Note, time.Now(), time.Now()).Scan(&rank.ID)
	if err != nil {
		errmsg("RankInsert QueryRow", err)
	}
	return rank.ID, err
}

// RankUpdate - save rank changes
func RankUpdate(rank Rank) error {
	_, err := pool.Exec(context.Background(), `
		UPDATE ranks SET
			name = $2,
			note = $3,
			updated_at = $4
		WHERE
			id = $1
	`, rank.ID, rank.Name, rank.Note, time.Now())
	if err != nil {
		errmsg("UpdateRank update", err)
	}
	return err
}

// RankDelete - delete rank by id
func RankDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.Exec(context.Background(), `
		DELETE FROM
			ranks
		WHERE
			id = $1
	`, id)
	if err != nil {
		errmsg("DeleteRank Exec", err)
	}
	return err
}

func rankCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			ranks (
				id bigserial primary key,
				name text,
				note text,
				created_at TIMESTAMP without time zone,
				updated_at TIMESTAMP without time zone default now(),
				UNIQUE (name)
			)
	`
	_, err := pool.Exec(context.Background(), str)
	if err != nil {
		errmsg("rankCreateTable exec", err)
	}
	return err
}
