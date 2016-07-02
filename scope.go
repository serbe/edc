package edc

import "log"

// Scope - struct for scope
type Scope struct {
	TableName struct{} `sql:"scopes"`
	ID        int64    `sql:"id" json:"id"`
	Name      string   `sql:"name" json:"name"`
	Notes     string   `sql:"notes, null" json:"notes"`
}

// GetScope - get one scope by id
func (e *EDc) GetScope(id int64) (scope Scope, err error) {
	if id == 0 {
		return
	}
	_, err = e.db.QueryOne(&scope, "SELECT * FROM scopes WHERE id = ? LIMIT 1", id)
	if err != nil {
		log.Println("GetScope e.db.QueryRow Scan ", err)
	}
	return
}

// GetScopeAll - get all scope
func (e *EDc) GetScopeAll() (scopes []Scope, err error) {
	_, err = e.db.Query(&scopes, "SELECT * FROM scopes")
	if err != nil {
		log.Println("GetScopeAll e.db.Query ", err)
		return
	}
	return
}

// CreateScope - create new scope
func (e *EDc) CreateScope(scope Scope) (err error) {
	err = e.db.Create(&scope)
	if err != nil {
		log.Println("CreateScope e.db.Exec ", err)
	}
	return
}

// UpdateScope - save scope changes
func (e *EDc) UpdateScope(scope Scope) (err error) {
	err = e.db.Update(&scope)
	if err != nil {
		log.Println("UpdateScope e.db.Exec ", err)
	}
	return
}

// DeleteScope - delete scope by id
func (e *EDc) DeleteScope(id int64) (err error) {
	if id == 0 {
		return
	}
	_, err = e.db.Exec("DELETE FROM scopes WHERE id = ?", id)
	if err != nil {
		log.Println("DeleteScope e.db.Exec ", err)
	}
	return
}

func (e *EDc) scopeCreateTable() (err error) {
	str := `CREATE TABLE IF NOT EXISTS scopes (id bigserial primary key, name text, notes text)`
	_, err = e.db.Exec(str)
	if err != nil {
		log.Println("scopeCreateTable e.db.Exec ", err)
	}
	return
}
