package edc

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
	err := pool.QueryRow(context.Background(), &kind).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetKind select", err)
	}
	return kind, err
}

// KindListGet - get kind for list by id
func KindListGet(id int64) (KindList, error) {
	var kind KindList
	err := pool.QueryRow(context.Background(), &Kind{}).
		Column("id", "name", "short_name", "note").
		Where("id = ?", id).
		Select(&kind)
	if err != nil {
		errmsg("GetKindList select", err)
	}
	return kind, err
}

// KindListAllGet - get all kind for list
func KindListAllGet() ([]KindList, error) {
	var kinds []KindList
	err := pool.QueryRow(context.Background(), &Kind{}).
		Column("id", "name", "short_name", "note").
		Order("name ASC").
		Select(&kinds)
	if err != nil {
		errmsg("GetKindListAll select", err)
	}
	return kinds, err
}

// KindSelectGet - get all kind for select
func KindSelectGet() ([]SelectItem, error) {
	var kinds []SelectItem
	err := pool.QueryRow(context.Background(), &Kind{}).
		Column("id", "name").
		Order("name ASC").
		Select(&kinds)
	if err != nil {
		errmsg("GetKindSelectAll select", err)
	}
	return kinds, err
}

// KindInsert - create new kind
func KindInsert(kind Kind) (int64, error) {
	err := pool.Insert(&kind)
	if err != nil {
		errmsg("CreateKind insert", err)
	}
	return kind.ID, nil
}

// KindUpdate - save kind changes
func KindUpdate(kind Kind) error {
	err := pool.Update(&kind)
	if err != nil {
		errmsg("UpdateKind update", err)
	}
	return err
}

// KindDelete - delete kind by id
func KindDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.QueryRow(context.Background(), &Kind{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteKind delete", err)
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
				updated_at
 TIMESTAMP without time zone default now(),
				UNIQUE(name)
			)
	`
	_, err := pool.Exec(str)
	if err != nil {
		errmsg("kindCreateTable exec", err)
	}
	return err
}
