package edc

// Tcc - struct for tcc
type Tcc struct {
	ID        int64   `sql:"id"              json:"id"         form:"id"         query:"id"`
	Address   string  `sql:"address,null"    json:"address"    form:"address"    query:"address"`
	ContactID int64   `sql:"contact_id,null" json:"contact_id" form:"contact_id" query:"contact_id"`
	Contact   Contact `sql:"-"               json:"contact"    form:"contact"    query:"contact"`
	CompanyID int64   `sql:"company_id,null" json:"company_id" form:"company_id" query:"company_id"`
	Company   Company `sql:"-"               json:"company"    form:"company"    query:"company"`
	Note      string  `sql:"note,null"       json:"note"       form:"note"       query:"note"`
	CreatedAt string  `sql:"created_at"      json:"-"`
	UpdatedAt string  `sql:"updated_at"      json:"-"`
}

// TccList - struct for tcc list
type TccList struct {
	ID        int64   `sql:"id"              json:"id"         form:"id"         query:"id"`
	Address   string  `sql:"address,null"    json:"address"    form:"address"    query:"address"`
	ContactID int64   `sql:"contact_id,null" json:"contact_id" form:"contact_id" query:"contact_id"`
	Contact   Contact `sql:"-"               json:"contact"    form:"contact"    query:"contact"`
	Note      string  `sql:"note,null"       json:"note"       form:"note"       query:"note"`
}

// GetTcc - get one tcc by id
func (e *Edb) GetTcc(id int64) (Tcc, error) {
	var tcc Tcc
	if id == 0 {
		return tcc, nil
	}
	err := e.db.Model(&tcc).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetTcc select", err)
	}
	return tcc, err
}

// GetTccList - get all tcc for list
func (e *Edb) GetTccList(id int64) (TccList, error) {
	var tccs TccList
	err := e.db.Model(&Tcc{}).
		Column("id", "address", "contact_id", "note").
		Where("id = ?", id).
		Select(&tccs)
	if err != nil {
		errmsg("GetTccList select", err)
	}
	return tccs, err
}

// GetTccListAll - get all tcc for list
func (e *Edb) GetTccListAll() ([]TccList, error) {
	var tccs []TccList
	err := e.db.Model(&Tcc{}).
		Column("id", "address", "contact_id", "note").
		Order("name ASC").
		Select(&tccs)
	if err != nil {
		errmsg("GetTccList select", err)
	}
	return tccs, err
}

// CreateTcc - create new tcc
func (e *Edb) CreateTcc(tcc Tcc) (int64, error) {
	err := e.db.Insert(&tcc)
	if err != nil {
		errmsg("CreateTcc insert", err)
	}
	return tcc.ID, err
}

// UpdateTcc - save tcc changes
func (e *Edb) UpdateTcc(tcc Tcc) error {
	err := e.db.Update(&tcc)
	if err != nil {
		errmsg("UpdateTcc update", err)
	}
	return err
}

// DeleteTcc - delete tcc by id
func (e *Edb) DeleteTcc(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Model(&Tcc{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteTcc delete", err)
	}
	return err
}

// func (e *Edb) tccCreateTable() error {
// 	str := `
// 		CREATE TABLE IF NOT EXISTS
// 			tccs (
// 				id         bigserial PRIMARY KEY,
// 				address    text,
// 				contact_id bigint,
// 				company_id bigint,
// 				note       text,
// 				created_at TIMESTAMP without time zone,
// 				updated_at TIMESTAMP without time zone default now(),
// 				UNIQUE(num_id, num_pass, type_id)
// 			)
// 	`
// 	_, err := e.db.Exec(str)
// 	if err != nil {
// 		errmsg("tccCreateTable exec", err)
// 	}
// 	return err
// }
