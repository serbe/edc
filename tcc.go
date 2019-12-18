package edc

import "context"

// Tcc - struct for tcc
type Tcc struct {
	ID        int64   `sql:"id"         json:"id"         form:"id"         query:"id"`
	Address   string  `sql:"address"    json:"address"    form:"address"    query:"address"`
	ContactID int64   `sql:"contact_id" json:"contact_id" form:"contact_id" query:"contact_id"`
	Contact   Contact `sql:"-"          json:"contact"    form:"contact"    query:"contact"`
	CompanyID int64   `sql:"company_id" json:"company_id" form:"company_id" query:"company_id"`
	Company   Company `sql:"-"          json:"company"    form:"company"    query:"company"`
	Note      string  `sql:"note"       json:"note"       form:"note"       query:"note"`
	CreatedAt string  `sql:"created_at" json:"-"`
	UpdatedAt string  `sql:"updated_at" json:"-"`
}

// TccList - struct for tcc list
type TccList struct {
	ID        int64   `sql:"id"         json:"id"         form:"id"         query:"id"`
	Address   string  `sql:"address"    json:"address"    form:"address"    query:"address"`
	ContactID int64   `sql:"contact_id" json:"contact_id" form:"contact_id" query:"contact_id"`
	Contact   Contact `sql:"-"          json:"contact"    form:"contact"    query:"contact"`
	Note      string  `sql:"note"       json:"note"       form:"note"       query:"note"`
}

// TccGet - get one tcc by id
func TccGet(id int64) (Tcc, error) {
	var tcc Tcc
	if id == 0 {
		return tcc, nil
	}
	err := pool.QueryRow(context.Background(), &tcc).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetTcc select", err)
	}
	return tcc, err
}

// TccListGet - get all tcc for list
func TccListGet() ([]TccList, error) {
	var tccs []TccList
	err := pool.QueryRow(context.Background(), &Tcc{}).
		Column("id", "address", "contact_id", "note").
		Order("name ASC").
		Select(&tccs)
	if err != nil {
		errmsg("GetTccList select", err)
	}
	return tccs, err
}

// TccInsert - create new tcc
func TccInsert(tcc Tcc) (int64, error) {
	err := pool.Insert(&tcc)
	if err != nil {
		errmsg("CreateTcc insert", err)
	}
	return tcc.ID, err
}

// TccUpdate - save tcc changes
func TccUpdate(tcc Tcc) error {
	err := pool.Update(&tcc)
	if err != nil {
		errmsg("UpdateTcc update", err)
	}
	return err
}

// TccDelete - delete tcc by id
func TccDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.QueryRow(context.Background(), &Tcc{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteTcc delete", err)
	}
	return err
}

// func tccCreateTable() error {
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
// 	_, err := pool.Exec(str)
// 	if err != nil {
// 		errmsg("tccCreateTable exec", err)
// 	}
// 	return err
// }
