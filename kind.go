package edc

import "log"

// Kind - struct for kind
type Kind struct {
	TableName struct{} `sql:"kinds"`
	ID        int64    `sql:"id" json:"id"`
	Name      string   `sql:"name" json:"name"`
	Notes     string   `sql:"notes, null" json:"notes"`
}

// GetKind - get one kind by id
func (e *EDc) GetKind(id int64) (kind Kind, err error) {
	if id == 0 {
		return
	}
	_, err = e.db.QueryOne(&kind, "SELECT * FROM kinds WHERE id = ? LIMIT 1", id)
	if err != nil {
		log.Println("GetKind e.db.QueryRow Scan ", err)
	}
	return
}

// GetKindAll - get all kinds
func (e *EDc) GetKindAll() (kinds []Kind, err error) {
	_, err = e.db.Query(&kinds, "SELECT * FROM kinds")
	if err != nil {
		log.Println("GetKindAll e.db.Query ", err)
		return
	}
	return
}

// CreateKind - create new kind
func (e *EDc) CreateKind(kind Kind) (err error) {
	err = e.db.Create(&kind)
	if err != nil {
		log.Println("CreateKind e.db.Create ", err)
	}
	return
}

// UpdateKind - save kind changes
func (e *EDc) UpdateKind(kind Kind) (err error) {
	err = e.db.Update(&kind)
	if err != nil {
		log.Println("UpdateKind e.db.Update ", err)
	}
	return
}

// DeleteKind - delete kind by id
func (e *EDc) DeleteKind(id int64) (err error) {
	if id == 0 {
		return
	}
	_, err = e.db.Exec("DELETE FROM kinds WHERE id = ?", id)
	if err != nil {
		log.Println("DeleteKind e.db.Exec ", err)
	}
	return
}

func (e *EDc) kindCreateTable() (err error) {
	str := `CREATE TABLE IF NOT EXISTS kinds (id bigserial primary key, name text, notes text)`
	_, err = e.db.Exec(str)
	if err != nil {
		log.Println("kindCreateTable e.db.Exec ", err)
	}
	return
}
