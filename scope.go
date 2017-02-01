package edc

// Scope - struct for scope
type Scope struct {
	ID        int64  `sql:"id" json:"id"`
	Name      string `sql:"name" json:"name"`
	Note      string `sql:"note, null" json:"note"`
	CreatedAt string `sql:"created_at" json:"created_at"`
	UpdatedAt string `sql:"updated_at" json:"updated_at"`
}

// GetScope - get one scope by id
func (e *Edb) GetScope(id int64) (Scope, error) {
	var scope Scope
	if id == 0 {
		return scope, nil
	}
	err := e.db.Model(&scope).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetScope select", err)
	}
	return scope, err
}

// GetScopeList - get all scope for list
func (e *Edb) GetScopeList() ([]Scope, error) {
	var scopes []Scope
	_, err := e.db.Query(&scopes, `SELECT id, name, note FROM scopes ORDER BY name ASC`)
	if err != nil {
		errmsg("GetScopeList query", err)
	}
	return scopes, err
}

// GetScopeSelect - get all scope for select
func (e *Edb) GetScopeSelect() ([]SelectItem, error) {
	var scopes []SelectItem
	_, err := e.db.Query(&scopes, `SELECT id, name FROM scopes ORDER BY name ASC`)
	if err != nil {
		errmsg("GetScopeSelect query", err)
	}
	return scopes, err
}

// CreateScope - create new scope
func (e *Edb) CreateScope(scope Scope) (int64, error) {
	err := e.db.Insert(&scope)
	if err != nil {
		errmsg("CreateScope insert", err)
	}
	return scope.ID, err
}

// UpdateScope - save scope changes
func (e *Edb) UpdateScope(scope Scope) error {
	err := e.db.Update(&scope)
	if err != nil {
		errmsg("UpdateScope update", err)
	}
	return err
}

// DeleteScope - delete scope by id
func (e *Edb) DeleteScope(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Model(&Scope{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteScope delete", err)
	}
	return err
}

func (e *Edb) scopeCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			scopes (
				id bigserial primary key,
				name text,
				note text,
				created_at TIMESTAMP without time zone,
				updated_at TIMESTAMP without time zone default now(),
				UNIQUE (name)
			)
	`
	_, err := e.db.Exec(str)
	if err != nil {
		errmsg("scopeCreateTable exec", err)
	}
	return err
}
