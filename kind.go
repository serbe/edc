package edc

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
	err := e.db.Model(&kind).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetKind select", err)
	}
	return kind, err
}

// GetKindList - get all kind for list
func (e *Edb) GetKindList() ([]Kind, error) {
	var kinds []Kind
	err := e.db.Model(&kinds).
		Column("id", "name", "note").
		Order("name ASC").
		Select()
	if err != nil {
		errmsg("GetKindList select", err)
	}
	return kinds, err
}

// GetKindSelect - get all kind for select
func (e *Edb) GetKindSelect() ([]SelectItem, error) {
	var kinds []SelectItem
	err := e.db.Model(&Kind{}).
		Column("id", "name").
		Order("name ASC").
		Select(&kinds)
	if err != nil {
		errmsg("GetKindSelect select", err)
	}
	return kinds, err
}

// CreateKind - create new kind
func (e *Edb) CreateKind(kind Kind) (int64, error) {
	err := e.db.Insert(&kind)
	if err != nil {
		errmsg("CreateKind insert", err)
	}
	return kind.ID, nil
}

// UpdateKind - save kind changes
func (e *Edb) UpdateKind(kind Kind) error {
	err := e.db.Update(&kind)
	if err != nil {
		errmsg("UpdateKind update", err)
	}
	return err
}

// DeleteKind - delete kind by id
func (e *Edb) DeleteKind(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Model(&Kind{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteKind delete", err)
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
		errmsg("kindCreateTable exec", err)
	}
	return err
}
