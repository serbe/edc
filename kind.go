package edc

import "log"

// Kind - struct for kind
type Kind struct {
	ID        int64  `sql:"id" json:"id"`
	Name      string `sql:"name" json:"name"`
	Note      string `sql:"note, null" json:"note"`
	CreatedAt string `sql:"created_at" json:"created_at"`
	UpdatedAt string `sql:"updated_at" json:"updated_at"`
}

// GetKind - get one kind by id
func (e *Edb) GetKind(id int64) (Kind, error) {
	var kind Kind
	if id == 0 {
		return kind, nil
	}
	err := e.db.Model(&kind).Where(`id = ?`, id).Select()
	if err != nil {
		log.Println("GetKind e.db.Select ", err)
		return kind, err
	}
	return kind, err
}

// GetKindList - get all kind for list
func (e *Edb) GetKindList() ([]Kind, error) {
	var kinds []Kind
	_, err := e.db.Query(&kinds, `
		SELECT
			id,
			name,
			note
		FROM
			kinds
		ORDER BY
			name ASC
	`)
	if err != nil {
		log.Println("GetKindList e.db.Query ", err)
		return []Kind{}, err
	}
	return kinds, err
}

// GetKindSelect - get all kind for select
func (e *Edb) GetKindSelect() ([]SelectItem, error) {
	var kinds []SelectItem
	_, err := e.db.Query(kinds, `
		SELECT
			id,
			name
		FROM
			kinds
		ORDER BY
			name ASC
	`)
	if err != nil {
		log.Println("GetKindSelect e.db.Query ", err)
		return []SelectItem{}, err
	}
	return kinds, err
}

// CreateKind - create new kind
func (e *Edb) CreateKind(kind Kind) (int64, error) {
	err := e.db.Insert(&kind)
	if err != nil {
		log.Println("CreateKind e.db.Create ", err)
		return 0, err
	}
	return kind.ID, nil
}

// UpdateKind - save kind changes
func (e *Edb) UpdateKind(kind Kind) error {
	err := e.db.Update(&kind)
	if err != nil {
		log.Println("UpdateKind e.db.Update ", err)
		return err
	}
	return err
}

// DeleteKind - delete kind by id
func (e *Edb) DeleteKind(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec(`
		DELETE FROM
			kinds
		WHERE
			id = ?
	`, id)
	if err != nil {
		log.Println("DeleteKind e.db.Exec ", id, err)
	}
	return err
}

func (e *Edb) kindCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			kinds (
				id bigserial primary key,
				name text,
				note text,
				created_at TIMESTAMP without time zone,
				updated_at TIMESTAMP without time zone default now(),
				UNIQUE(name)
			)
	`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("kindCreateTable e.db.Exec ", err)
	}
	return err
}
