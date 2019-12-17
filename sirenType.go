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

// SirenTypeGet - get one sirenType by id
func SirenTypeGet(id int64) (SirenType, error) {
	var sirenType SirenType
	if id == 0 {
		return sirenType, nil
	}
	err := pool.QueryRow(context.Background(), &sirenType).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetSirenType select", err)
	}
	return sirenType, err
}

// SirenTypeListGet - get sirenType for list by id
func SirenTypeListGet(id int64) (SirenTypeList, error) {
	var sirenType SirenTypeList
	err := pool.QueryRow(context.Background(), &SirenType{}).
		Column("id", "name", "radius", "note").
		Where("id = ?", id).
		Select(&sirenType)
	if err != nil {
		errmsg("GetSirenTypeList select", err)
	}
	return sirenType, err
}

// SirenTypeListAllGet - get all sirenType for list
func SirenTypeListAllGet() ([]SirenTypeList, error) {
	var sirenTypes []SirenTypeList
	err := pool.QueryRow(context.Background(), &SirenType{}).
		Column("id", "name", "radius", "note").
		Order("name ASC").
		Select(&sirenTypes)
	if err != nil {
		errmsg("GetSirenTypeList select", err)
	}
	return sirenTypes, err
}

// SirenTypeSelectGet - get sirenType for select by id
func SirenTypeSelectGet(id int64) ([]SelectItem, error) {
	var sirenTypes []SelectItem
	err := pool.QueryRow(context.Background(), &SirenType{}).
		Column("id", "name").
		Where("id = ?", id).
		Select(&sirenTypes)
	if err != nil {
		errmsg("GetSirenTypeSelect Select", err)
	}
	return sirenTypes, err
}

// SirenTypeSelectGet - get all sirenType for select
func SirenTypeSelectGet() ([]SelectItem, error) {
	var sirenTypes []SelectItem
	err := pool.QueryRow(context.Background(), &SirenType{}).
		Column("id", "name").
		Order("name ASC").
		Select(&sirenTypes)
	if err != nil {
		errmsg("GetSirenTypeSelect Select", err)
	}
	return sirenTypes, err
}

// SirenTypeInsert - create new sirenType
func SirenTypeInsert(sirenType SirenType) (int64, error) {
	err := pool.Insert(&sirenType)
	if err != nil {
		errmsg("CreateSirenType insert", err)
	}
	return sirenType.ID, nil
}

// SirenTypeUpdate - save sirenType changes
func SirenTypeUpdate(sirenType SirenType) error {
	err := pool.Update(&sirenType)
	if err != nil {
		errmsg("UpdateSirenType update", err)
	}
	return err
}

// SirenTypeDelete - delete sirenType by id
func SirenTypeDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.QueryRow(context.Background(), &SirenType{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteSirenTypedelete", err)
	}
	return err
}

func sirenTypeCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			siren_types (
				id         bigserial primary key,
				name       text,
				radius     bigint,
				note       text,
				created_at TIMESTAMP without time zone,
				updated_at
 TIMESTAMP without time zone default now(),
				UNIQUE(name, radius)
			);`
	_, err := pool.Exec(str)
	if err != nil {
		errmsg("sirenCreateTable exec", err)
	}
	return err
}
