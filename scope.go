package edc

import "context"

// Scope - struct for scope
type Scope struct {
	ID        int64  `sql:"id"         json:"id"   form:"id"   query:"id"`
	Name      string `sql:"name"       json:"name" form:"name" query:"name"`
	Note      string `sql:"note"       json:"note" form:"note" query:"note"`
	CreatedAt string `sql:"created_at" json:"-"`
	UpdatedAt string `sql:"updated_at" json:"-"`
}

// ScopeList - struct for scope list
type ScopeList struct {
	ID   int64  `sql:"id"   json:"id"   form:"id"   query:"id"`
	Name string `sql:"name" json:"name" form:"name" query:"name"`
	Note string `sql:"note" json:"note" form:"note" query:"note"`
}

// ScopeGet - get one scope by id
func ScopeGet(id int64) (Scope, error) {
	var scope Scope
	if id == 0 {
		return scope, nil
	}
	err := pool.QueryRow(context.Background(), &scope).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetScope select", err)
	}
	return scope, err
}

// ScopeListGet - get all scope for list
func ScopeListGet() ([]ScopeList, error) {
	var scopes []ScopeList
	err := pool.QueryRow(context.Background(), &Scope{}).
		Column("id", "name", "note").
		Order("name ASC").
		Select(&scopes)
	if err != nil {
		errmsg("GetScopeListAll select", err)
	}
	return scopes, err
}

// ScopeSelectGet - get all scope for select
func ScopeSelectGet() ([]SelectItem, error) {
	var scopes []SelectItem
	err := pool.QueryRow(context.Background(), &Scope{}).
		Column("id", "name").
		Order("name ASC").
		Select(&scopes)
	if err != nil {
		errmsg("GetScopeSelectAll query", err)
	}
	rows, err := pool.Query(context.Background(), `
		SELECT
			id,
			name
		FROM
			companies
		ORDER BY
			name ASC
	`)
	if err != nil {
		errmsg("CompanySelectGet Query", err)
	}
	for rows.Next() {
		var company SelectItem
		err := rows.Scan(&company.ID, &company.Name)
		if err != nil {
			errmsg("CompanySelectGet select", err)
			return companies, err
		}
		companies = append(companies, company)
	}
	return companies, rows.Err()
	return scopes, err
}

// ScopeInsert - create new scope
func ScopeInsert(scope Scope) (int64, error) {
	err := pool.Insert(&scope)
	if err != nil {
		errmsg("CreateScope insert", err)
	}
	return scope.ID, err
}

// ScopeUpdate - save scope changes
func ScopeUpdate(scope Scope) error {
	err := pool.Update(&scope)
	if err != nil {
		errmsg("UpdateScope update", err)
	}
	return err
}

// ScopeDelete - delete scope by id
func ScopeDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.QueryRow(context.Background(), &Scope{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteScope delete", err)
	}
	return err
}

func scopeCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			scopes (
				id bigserial primary key,
				name text,
				note text,
				created_at TIMESTAMP without time zone,
				updated_at
 TIMESTAMP without time zone default now(),
				UNIQUE (name)
			)
	`
	_, err := pool.Exec(str)
	if err != nil {
		errmsg("scopeCreateTable exec", err)
	}
	return err
}
