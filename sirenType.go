package edc

// SirenType - struct for sirenType
type SirenType struct {
	ID        int64  `sql:"id"         json:"id"            form:"id"     query:"id"`
	Name      string `sql:"name"       json:"name"          form:"name"   query:"name"`
	Radius    int64  `sql:"radius"     json:"radius,string" form:"radius" query:"radius"`
	Note      string `sql:"note"       json:"note"          form:"note"   query:"note"`
	CreatedAt string `sql:"created_at" json:"-"`
	UpdatedAt string `sql:"updated_at" json:"-"`
}

// SirenTypeList - struct for sirenType list
type SirenTypeList struct {
	ID     int64  `sql:"id"     json:"id"            form:"id"     query:"id"`
	Name   string `sql:"name"   json:"name"          form:"name"   query:"name"`
	Radius int64  `sql:"radius" json:"radius,string" form:"radius" query:"radius"`
	Note   string `sql:"note"   json:"note"          form:"note"   query:"note"`
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

// GetSirenTypeList - get sirenType for list by id
func (e *Edb) GetSirenTypeList(id int64) (SirenTypeList, error) {
	var sirenType SirenTypeList
	err := e.db.Model(&SirenType{}).
		Column("id", "name", "radius", "note").
		Where("id = ?", id).
		Select(&sirenType)
	if err != nil {
		errmsg("GetSirenTypeList select", err)
	}
	return sirenType, err
}

// GetSirenTypeListAll - get all sirenType for list
func (e *Edb) GetSirenTypeListAll() ([]SirenTypeList, error) {
	var sirenTypes []SirenTypeList
	err := e.db.Model(&SirenType{}).
		Column("id", "name", "radius", "note").
		Order("name ASC").
		Select(&sirenTypes)
	if err != nil {
		errmsg("GetSirenTypeList select", err)
	}
	return sirenTypes, err
}

// GetSirenTypeSelect - get sirenType for select by id
func (e *Edb) GetSirenTypeSelect(id int64) ([]SelectItem, error) {
	var sirenTypes []SelectItem
	err := e.db.Model(&SirenType{}).
		Column("id", "name").
		Where("id = ?", id).
		Select(&sirenTypes)
	if err != nil {
		errmsg("GetSirenTypeSelect Select", err)
	}
	return sirenTypes, err
}

// GetSirenTypeSelectAll - get all sirenType for select
func (e *Edb) GetSirenTypeSelectAll() ([]SelectItem, error) {
	var sirenTypes []SelectItem
	err := e.db.Model(&SirenType{}).
		Column("id", "name").
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
			siren_types (
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
