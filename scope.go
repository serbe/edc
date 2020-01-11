package edc

import "context"

import "time"

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
	scope.ID = id
	err := pool.QueryRow(context.Background(), `
		SELECT
			name,
			note,
			created_at,
			updated_at
		FROM
			scopes
		WHERE
			id = $1
	`, id).Scan(scope.Name, scope.Note, scope.CreatedAt, scope.UpdatedAt)
	if err != nil {
		errmsg("ScopeGet QueryRow", err)
	}
	return scope, err
}

// ScopeListGet - get all scope for list
func ScopeListGet() ([]ScopeList, error) {
	var scopes []ScopeList
	rows, err := pool.Query(context.Background(), `
		SELECT
			id,
			name,
			note
		FROM
			scopes
		ORDER BY
			name ASC
	`)
	if err != nil {
		errmsg("ScopeListGet Query", err)
	}
	for rows.Next() {
		var scope ScopeList
		err := rows.Scan(&scope.ID, &scope.Name, &scope.Note)
		if err != nil {
			errmsg("ScopeListGet Scan", err)
			return scopes, err
		}
		scopes = append(scopes, scope)
	}
	return scopes, rows.Err()
}

// ScopeSelectGet - get all scope for select
func ScopeSelectGet() ([]SelectItem, error) {
	var scopes []SelectItem
	rows, err := pool.Query(context.Background(), `
		SELECT
			id,
			name
		FROM
			scopes
		ORDER BY
			name ASC
	`)
	if err != nil {
		errmsg("CompanySelectGet Query", err)
	}
	for rows.Next() {
		var scope SelectItem
		err := rows.Scan(&scope.ID, &scope.Name)
		if err != nil {
			errmsg("ScopeSelectGet Scan", err)
			return scopes, err
		}
		scopes = append(scopes, scope)
	}
	return scopes, rows.Err()
}

// ScopeInsert - create new scope
func ScopeInsert(scope Scope) (int64, error) {
	err := pool.QueryRow(context.Background(), `
		INSERT INTO scopes
		(
			name,
			note,
			created_at,
			updated_at
		)
		VALUES
		(
			$1,
			$2,
			$3,
			$4
		)
		RETURNING
			id
	`, scope.Name, scope.Note, time.Now(), time.Now()).Scan(&scope.ID)
	if err != nil {
		errmsg("ScopeInsert QueryRow", err)
	}
	return scope.ID, err
}

// ScopeUpdate - save scope changes
func ScopeUpdate(scope Scope) error {
	_, err := pool.Exec(context.Background(), `
		UPDATE scopes SET
			name = $2,
			note = $3,
			updated_at = $4
		WHERE
			id = $1
	`, scope.ID, scope.Name, scope.Note, time.Now())
	if err != nil {
		errmsg("ScopeUpdate Exec", err)
	}
	return err
}

// ScopeDelete - delete scope by id
func ScopeDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.Exec(context.Background(), `
		DELETE FROM
			scopes
		WHERE
			id = $1
	`, id)
	if err != nil {
		errmsg("ScopeDelete Exec", err)
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
				updated_at TIMESTAMP without time zone default now(),
				UNIQUE (name)
			)
	`
	_, err := pool.Exec(context.Background(), str)
	if err != nil {
		errmsg("scopeCreateTable Exec", err)
	}
	return err
}
