package edc

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
	err := e.db.Model(&rank).
		Where("id = ?", id).
		Select()
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
	var ranks []SelectItem
	_, err := e.db.Query(&ranks, `
		SELECT
			id,
			name
		FROM
			ranks
		ORDER BY
			name ASC
	`)
	if err != nil {
		errmsg("GetRankSelect query", err)
	}
	return ranks, err
}

// CreateRank - create new rank
func (e *Edb) CreateRank(rank Rank) (int64, error) {
	err := e.db.Insert(&rank)
	if err != nil {
		errmsg("CreateRank insert", err)
	}
	return rank.ID, err
}

// UpdateRank - save rank changes
func (e *Edb) UpdateRank(rank Rank) error {
	err := e.db.Update(&rank)
	if err != nil {
		errmsg("UpdateRank update", err)
	}
	return err
}

// DeleteRank - delete rank by id
func (e *Edb) DeleteRank(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Model(&Rank{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteRank delete", err)
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
				updated_at TIMESTAMP without time zone default now(),
				UNIQUE (name)
			)
	`
	_, err := e.db.Exec(str)
	if err != nil {
		errmsg("rankCreateTable exec", err)
	}
	return err
}
