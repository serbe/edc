package edc

// HideoutType - struct for hideoutType
type HideoutType struct {
	ID        int64  `sql:"id"         json:"id"   form:"id"   query:"id"`
	Name      string `sql:"name"       json:"name" form:"name" query:"name"`
	Note      string `sql:"note"       json:"note" form:"note" query:"note"`
	CreatedAt string `sql:"created_at" json:"-"`
	UpdatedAt string `sql:"updated_at" json:"-"`
}

// HideoutTypeList - struct for hideoutType list
type HideoutTypeList struct {
	ID   int64  `sql:"id"   json:"id"   form:"id"   query:"id"`
	Name string `sql:"name" json:"name" form:"name" query:"name"`
	Note string `sql:"note" json:"note" form:"note" query:"note"`
}

// GetHideoutType - get one hideoutType by id
func (e *Edb) GetHideoutType(id int64) (HideoutType, error) {
	var hideoutType HideoutType
	if id == 0 {
		return hideoutType, nil
	}
	err := e.db.Model(&hideoutType).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetHideoutType select", err)
	}
	return hideoutType, err
}

// GetHideoutTypeList - get hideoutType for list by id
func (e *Edb) GetHideoutTypeList(id int64) (HideoutTypeList, error) {
	var hideoutType HideoutTypeList
	err := e.db.Model(&HideoutType{}).
		Column("id", "name", "note").
		Where("id = ?", id).
		Select(&hideoutType)
	if err != nil {
		errmsg("GetHideoutTypeList select", err)
	}
	return hideoutType, err
}

// GetHideoutTypeListAll - get all hideoutType for list
func (e *Edb) GetHideoutTypeListAll() ([]HideoutTypeList, error) {
	var hideoutTypes []HideoutTypeList
	err := e.db.Model(&HideoutType{}).
		Column("id", "name", "note").
		Order("name ASC").
		Select(&hideoutTypes)
	if err != nil {
		errmsg("GetHideoutTypeList select", err)
	}
	return hideoutTypes, err
}

// GetHideoutTypeSelect - get hideoutType for select by id
func (e *Edb) GetHideoutTypeSelect(id int64) ([]SelectItem, error) {
	var hideoutTypes []SelectItem
	err := e.db.Model(&HideoutType{}).
		Column("id", "name").
		Where("id = ?", id).
		Select(&hideoutTypes)
	if err != nil {
		errmsg("GetHideoutTypeSelect Select", err)
	}
	return hideoutTypes, err
}

// GetHideoutTypeSelectAll - get all hideoutType for select
func (e *Edb) GetHideoutTypeSelectAll() ([]SelectItem, error) {
	var hideoutTypes []SelectItem
	err := e.db.Model(&HideoutType{}).
		Column("id", "name").
		Order("name ASC").
		Select(&hideoutTypes)
	if err != nil {
		errmsg("GetHideoutTypeSelect Select", err)
	}
	return hideoutTypes, err
}

// CreateHideoutType - create new hideoutType
func (e *Edb) CreateHideoutType(hideoutType HideoutType) (int64, error) {
	err := e.db.Insert(&hideoutType)
	if err != nil {
		errmsg("CreateHideoutType insert", err)
	}
	return hideoutType.ID, nil
}

// UpdateHideoutType - save hideoutType changes
func (e *Edb) UpdateHideoutType(hideoutType HideoutType) error {
	err := e.db.Update(&hideoutType)
	if err != nil {
		errmsg("UpdateHideoutType update", err)
	}
	return err
}

// DeleteHideoutType - delete hideoutType by id
func (e *Edb) DeleteHideoutType(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Model(&HideoutType{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteHideoutTypedelete", err)
	}
	return err
}

func (e *Edb) hideoutTypeCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			hideout_types (
				id         bigserial primary key,
				name       text,
				note       text,
				created_at TIMESTAMP without time zone,
				updated_at TIMESTAMP without time zone default now(),
				UNIQUE(name)
			);`
	_, err := e.db.Exec(str)
	if err != nil {
		errmsg("hideoutCreateTable exec", err)
	}
	return err
}
