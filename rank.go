package edc

import "fmt"

// Rank - struct for rank
type Rank struct {
	ID    int64  `sql:"id" json:"id"`
	Name  string `sql:"name" json:"name"`
	Notes string `sql:"notes, null" json:"notes"`
}

// GetRank - get one rank dy id
func (e *EDc) GetRank(id int64) (rank Rank, err error) {
	if id == 0 {
		return rank, nil
	}
	_, err = e.db.QueryOne(&rank, "SELECT * FROM ranks WHERE id = ? LIMIT 1", id)
	if err != nil {
		return rank, fmt.Errorf("GetRank e.db.QueryRow Scan: %s", err)
	}
	return
}

// GetRankAll - get all rank
func (e *EDc) GetRankAll() (ranks []Rank, err error) {
	_, err = e.db.Query(&ranks, "SELECT * FROM ranks")
	if err != nil {
		return ranks, fmt.Errorf("GetRankAll e.db.Query: %s", err)
	}
	return
}

// CreateRank - create new rank
func (e *EDc) CreateRank(rank Rank) (err error) {
	err = e.db.Create(&rank)
	if err != nil {
		return fmt.Errorf("CreateRank e.db.Exec: %s", err)
	}
	return
}

// UpdateRank - save rank changes
func (e *EDc) UpdateRank(rank Rank) (err error) {
	err = e.db.Update(&rank)
	if err != nil {
		return fmt.Errorf("UpdateRank e.db.Exec: %s", err)
	}
	return
}

// DeleteRank - delete rank by id
func (e *EDc) DeleteRank(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec("DELETE FROM ranks WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("DeleteRank e.db.Exec: %s", err)
	}
	return nil
}

func (e *EDc) rankCreateTable() (err error) {
	str := `CREATE TABLE IF NOT EXISTS ranks (id bigserial primary key, name text, notes text)`
	_, err = e.db.Exec(str)
	if err != nil {
		return fmt.Errorf("rankCreateTable e.db.Exec: %s", err)
	}
	return
}
