package edc

// SirenType - struct for sirenType
type SirenType struct {
	ID        int64  `sql:"id"         json:"id"`
	Name      string `sql:"name"       json:"name"`
	Radius    int64  `sql:"radius"     json:"radius"`
	Note      string `sql:"note, null" json:"note"`
	CreatedAt string `sql:"created_at" json:"-"`
	UpdatedAt string `sql:"updated_at" json:"-"`
}

// GetSirenType - get one sirenType by id
func (e *Edb) GetSirenType(id int64) (SirenType, error) {
	var sirenType SirenType
	if id == 0 {
		return sirenType, nil
	}
	err := e.db.Model(&sirenType).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetSirenType select", err)
	}
	return sirenType, err
}

// GetSirenTypeList - get all sirenType for list
func (e *Edb) GetSirenTypeList() ([]SirenType, error) {
	var sirenTypes []SirenType
	err := e.db.Model(&sirenTypes).
		Order("name ASC").
		Select()
	if err != nil {
		errmsg("GetSirenTypeList select", err)
	}
	return sirenTypes, err
}

// GetSirenTypeSelect - get all sirenType for select
func (e *Edb) GetSirenTypeSelect() ([]SelectItem, error) {
	var sirenTypes []SelectItem
	err := e.db.Model(&SirenType{}).
		Order("name ASC").
		Select(&sirenTypes)
	if err != nil {
		errmsg("GetSirenTypeSelect Select", err)
	}
	return sirenTypes, err
}

// CreateSirenType - create new sirenType
func (e *Edb) CreateSirenType(sirenType SirenType) (int64, error) {
	err := e.db.Insert(&sirenType)
	if err != nil {
		errmsg("CreateSirenType insert", err)
	}
	return sirenType.ID, nil
}

// UpdateSirenType - save sirenType changes
func (e *Edb) UpdateSirenType(sirenType SirenType) error {
	err := e.db.Update(&sirenType)
	if err != nil {
		errmsg("UpdateSirenType update", err)
	}
	return err
}

// DeleteSirenType - delete sirenType by id
func (e *Edb) DeleteSirenType(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Model(&SirenType{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteSirenTypedelete", err)
	}
	return err
}

func (e *Edb) sirenTypeCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			sirenTypes (
				id         bigserial primary key,
				name       text,
				radius     bigint,
				note       text,
				created_at TIMESTAMP without time zone,
				updated_at TIMESTAMP without time zone default now(),
				UNIQUE(name, radius)
			);`
	_, err := e.db.Exec(str)
	if err != nil {
		errmsg("sirenCreateTable exec", err)
	}
	return err
}
