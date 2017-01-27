package edc

import "log"

// Rank - struct for rank
type Rank struct {
	ID        int64  `sql:"id" json:"id"`
	Name      string `sql:"name" json:"name"`
	Note      string `sql:"note, null" json:"note"`
	CreatedAt string `sql:"created_at" json:"created_at"`
	UpdatedAt string `sql:"updated_at" json:"updated_at"`
}

// GetRank - get one rank by id
func (e *Edb) GetRank(id int64) (Rank, error) {
	var rank Rank
	if id == 0 {
		return rank, nil
	}
	err := e.db.Model(&rank).Where("id = ?", id).Select()
	if err != nil {
		errmsg("GetRank select", err)
	}
	return rank, err
}

// GetRankList - get all rank for list
func (e *Edb) GetRankList() ([]Rank, error) {
	var ranks []Rank
	_, err := e.db.Query(&ranks, `
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
		errmsg("GetRankList query", err)
	}
	return ranks, err
}

// GetRankSelect - get all rank for select
func (e *Edb) GetRankSelect() ([]SelectItem, error) {
	rows, err := e.db.Query(`
		SELECT
			id,
			name
		FROM
			ranks
		ORDER BY
			name ASC
	`)
	if err != nil {
		log.Println("GetRankSelect e.db.Query ", err)
		return []SelectItem{}, err
	}
	ranks, err := scanRanksSelect(rows)
	return ranks, err
}

// CreateRank - create new rank
func (e *Edb) CreateRank(rank Rank) (int64, error) {
	stmt, err := e.db.Prepare(`
		INSERT INTO
			ranks (
				name,
				note,
				created_at
			) VALUES (
				$1,
				$2,
				now()
			)
		RETURNING
			id
	`)
	if err != nil {
		log.Println("CreateRank e.db.Prepare ", err)
		return 0, err
	}
	err = stmt.QueryRow(s2n(rank.Name), s2n(rank.Note)).Scan(&rank.ID)
	if err != nil {
		log.Println("CreateRank db.QueryRow ", err)
	}
	return rank.ID, err
}

// UpdateRank - save rank changes
func (e *Edb) UpdateRank(s Rank) error {
	stmt, err := e.db.Prepare(`
		UPDATE
			ranks
		SET
			name = $2,
			note = $3,
			updated_at = now()
		WHERE
			id = $1
	`)
	if err != nil {
		log.Println("UpdateRank e.db.Prepare ", err)
		return err
	}
	_, err = stmt.Exec(i2n(s.ID), s2n(s.Name), s2n(s.Note))
	if err != nil {
		log.Println("UpdateRank stmt.Exec ", err)
	}
	return err
}

// DeleteRank - delete rank by id
func (e *Edb) DeleteRank(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec(`
		DELETE FROM
			ranks
		WHERE
			id = $1
	`, id)
	if err != nil {
		log.Println("DeleteRank e.db.Exec ", id, err)
	}
	return err
}

func (e *Edb) rankCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			ranks (
				id bigserial primary key,
				name text,
				note text,
				created_at TIMESTAMP without time zone,
				updated_at TIMESTAMP without time zone,
				UNIQUE (name)
			)
	`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("rankCreateTable e.db.Exec ", err)
	}
	return err
}
