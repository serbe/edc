package edc

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
	err := pool.QueryRow(context.Background(), &rank).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetRank select", err)
	}
	return rank, err
}

// RankListGet - get rank for list by id
func RankListGet(id int64) (RankList, error) {
	var rank RankList
	err := pool.QueryRow(context.Background(), &Rank{}).
		Column("id", "name", "note").
		Where("id = ?", id).
		Select(&rank)
	if err != nil {
		errmsg("GetRankList query", err)
	}
	return rank, err
}

// RankListAllGet - get all rank for list
func RankListAllGet() ([]RankList, error) {
	var ranks []RankList
	err := pool.QueryRow(context.Background(), &Rank{}).
		Column("id", "name", "note").
		Order("name ASC").
		Select(&ranks)
	if err != nil {
		errmsg("GetRankListAll query", err)
	}
	return ranks, err
}

// RankSelectGet - get all rank for select
func RankSelectGet(id int64) (SelectItem, error) {
	var rank SelectItem
	if id == 0 {
		return rank, nil
	}
	err := pool.QueryRow(context.Background(), &Rank{}).
		Column("id", "name").
		Where("id = ?", id).
		Order("name ASC").
		Select(&rank)
	if err != nil {
		errmsg("GetRankSelect query", err)
	}
	return rank, err
}

// RankSelectGet - get all rank for select
func RankSelectGet() ([]SelectItem, error) {
	var ranks []SelectItem
	err := pool.QueryRow(context.Background(), &Rank{}).
		Column("id", "name").
		Order("name ASC").
		Select(&ranks)
	if err != nil {
		errmsg("GetRankSelectAll query", err)
	}
	return ranks, err
}

// RankInsert - create new rank
func RankInsert(rank Rank) (int64, error) {
	err := pool.Insert(&rank)
	if err != nil {
		errmsg("CreateRank insert", err)
	}
	return rank.ID, err
}

// RankUpdate - save rank changes
func RankUpdate(rank Rank) error {
	err := pool.Update(&rank)
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
	_, err := pool.QueryRow(context.Background(), &Rank{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteRank delete", err)
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
				updated_at
 TIMESTAMP without time zone default now(),
				UNIQUE (name)
			)
	`
	_, err := pool.Exec(str)
	if err != nil {
		errmsg("rankCreateTable exec", err)
	}
	return err
}
