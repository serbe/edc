package edc

// Scope - struct for scope
type Scope struct {
	ID        int64  `sql:"id"         json:"id"   form:"id"   query:"id"`
	Name      string `sql:"name"       json:"name" form:"name" query:"name"`
	Note      string `sql:"note,null"  json:"note" form:"note" query:"note"`
	CreatedAt string `sql:"created_at" json:"-"`
	UpdatedAt string `sql:"updated_at" json:"-"`
}

// ScopeList - struct for scope list
type ScopeList struct {
	ID   int64  `sql:"id"        json:"id"   form:"id"   query:"id"`
	Name string `sql:"name"      json:"name" form:"name" query:"name"`
	Note string `sql:"note,null" json:"note" form:"note" query:"note"`
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

// GetScopeList - get scope for list by id
func (e *Edb) GetScopeList(id int64) (ScopeList, error) {
	var scope ScopeList
	err := e.db.Model(&Scope{}).
		Column("id", "name", "note").
		Where("id = ?", id).
		Select(&scope)
	if err != nil {
		errmsg("GetScopeList select", err)
	}
	return scope, err
}

// GetScopeListAll - get all scope for list
func (e *Edb) GetScopeListAll() ([]ScopeList, error) {
	var scopes []ScopeList
	err := e.db.Model(&Scope{}).
		Column("id", "name", "note").
		Order("name ASC").
		Select(&scopes)
	if err != nil {
		errmsg("GetScopeListAll select", err)
	}
	return scopes, err
}

// GetScopeSelect - get scope for select
func (e *Edb) GetScopeSelect(id int64) (SelectItem, error) {
	var scope SelectItem
	if id == 0 {
		return scope, nil
	}
	err := e.db.Model(&Scope{}).
		Column("id", "name").
		Where("id = ?", id).
		Select(&scope)
	if err != nil {
		errmsg("GetScopeSelect select", err)
	}
	return scope, err
}

// GetScopeSelectAll - get all scope for select
func (e *Edb) GetScopeSelectAll() ([]SelectItem, error) {
	var scopes []SelectItem
	err := e.db.Model(&Scope{}).
		Column("id", "name").
		Order("name ASC").
		Select(&scopes)
	if err != nil {
		errmsg("GetScopeSelectAll query", err)
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
