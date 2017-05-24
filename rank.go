package edc

// Rank - struct for rank
type Rank struct {
	ID        int64  `sql:"id"         json:"id" form:"id" query:"id"`
	Name      string `sql:"name"       json:"name" form:"name" query:"name"`
	Note      string `sql:"note,null"  json:"note" form:"note" query:"note"`
	CreatedAt string `sql:"created_at" json:"-" form:"-" query:"-"`
	UpdatedAt string `sql:"updated_at" json:"-" form:"-" query:"-"`
}

// RankList - struct for rank list
type RankList struct {
	ID   int64  `sql:"id"        json:"id" form:"id" query:"id"`
	Name string `sql:"name"      json:"name" form:"name" query:"name"`
	Note string `sql:"note,null" json:"note" form:"note" query:"note"`
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

// GetRankList - get rank for list by id
func (e *Edb) GetRankList(id int64) (RankList, error) {
	var rank RankList
	err := e.db.Model(&Rank{}).
		Column("id", "name", "note").
		Where("id = ?", id).
		Select(&rank)
	if err != nil {
		errmsg("GetRankList query", err)
	}
	return rank, err
}

// GetRankListAll - get all rank for list
func (e *Edb) GetRankListAll() ([]RankList, error) {
	var ranks []RankList
	err := e.db.Model(&Rank{}).
		Column("id", "name", "note").
		Order("name ASC").
		Select(&ranks)
	if err != nil {
		errmsg("GetRankListAll query", err)
	}
	return ranks, err
}

// GetRankSelect - get all rank for select
func (e *Edb) GetRankSelect(id int64) (SelectItem, error) {
	var rank SelectItem
	if id == 0 {
		return rank, nil
	}
	err := e.db.Model(&Rank{}).
		Column("id", "name").
		Where("id = ?", id).
		Order("name ASC").
		Select(&rank)
	if err != nil {
		errmsg("GetRankSelect query", err)
	}
	return rank, err
}

// GetRankSelectAll - get all rank for select
func (e *Edb) GetRankSelectAll() ([]SelectItem, error) {
	var ranks []SelectItem
	err := e.db.Model(&Rank{}).
		Column("id", "name").
		Order("name ASC").
		Select(&ranks)
	if err != nil {
		errmsg("GetRankSelectAll query", err)
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
