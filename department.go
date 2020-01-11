package edc

import (
	"context"
	"time"
)

// Department - struct for department
type Department struct {
	ID        int64  `sql:"id"         json:"id"   form:"id"   query:"id"`
	Name      string `sql:"name"       json:"name" form:"name" query:"name"`
	Note      string `sql:"note"       json:"note" form:"note" query:"note"`
	CreatedAt string `sql:"created_at" json:"-"    form:"-"    query:"-"`
	UpdatedAt string `sql:"updated_at" json:"-"    form:"-"    query:"-"`
}

// DepartmentList - struct for list of departments
type DepartmentList struct {
	ID   int64  `sql:"id"   json:"id"   form:"id"   query:"id"`
	Name string `sql:"name" json:"name" form:"name" query:"name"`
	Note string `sql:"note" json:"note" form:"note" query:"note"`
}

// DepartmentGet - get one department by id
func DepartmentGet(id int64) (Department, error) {
	var department Department
	if id == 0 {
		return department, nil
	}
	department.ID = id
	err := pool.QueryRow(context.Background(), `
		SELECT
			name,
			note,
			created_at,
			updated_at
		FROM
			departments
		WHERE
			id = $1
	`, id).Scan(&department.Name, &department.Note, time.Now(), time.Now())
	if err != nil {
		errmsg("DepartmentGet QueryRow", err)
	}
	return department, err
}

// DepartmentListGet - get all department for list
func DepartmentListGet() ([]DepartmentList, error) {
	var departments []DepartmentList
	rows, err := pool.Query(context.Background(), `
		SELECT
			id,
			name,
			note
		FROM
			departments
		ORDER BY
			name ASC
	`)
	if err != nil {
		errmsg("DepartmentListGet Query", err)
	}
	for rows.Next() {
		var department DepartmentList
		err := rows.Scan(&department.ID, &department.Name, &department.Note)
		if err != nil {
			errmsg("DepartmentListGet Scan", err)
			return departments, err
		}
		departments = append(departments, department)
	}
	return departments, rows.Err()
}

// DepartmentSelectGet - get all department for select
func DepartmentSelectGet() ([]SelectItem, error) {
	var departments []SelectItem
	rows, err := pool.Query(context.Background(), `
		SELECT
			id,
			name
		FROM
			departments
		ORDER BY
			name ASC
	`)
	if err != nil {
		errmsg("CompanySelectGet Query", err)
	}
	for rows.Next() {
		var department SelectItem
		err := rows.Scan(&department.ID, &department.Name)
		if err != nil {
			errmsg("CompanySelectGet Scan", err)
			return departments, err
		}
		departments = append(departments, department)
	}
	return departments, rows.Err()
}

// DepartmentInsert - create new department
func DepartmentInsert(department Department) (int64, error) {
	err := pool.QueryRow(context.Background(), `
		INSERT INTO departments
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
	`, department.Name, department.Note, time.Now(), time.Now()).Scan(&department.ID)
	if err != nil {
		errmsg("DepartmentInsert QueryRow", err)
	}
	return department.ID, nil
}

// DepartmentUpdate - save department changes
func DepartmentUpdate(department Department) error {
	_, err := pool.Exec(context.Background(), `
		UPDATE departments SET
			name = $2,
			note = $3,
			updated_at = $4
		WHERE
			id = $1
	`, department.ID, department.Name, department.Note, time.Now())
	if err != nil {
		errmsg("DepartmentUpdate Exec", err)
	}
	return err
}

// DepartmentDelete - delete department by id
func DepartmentDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.Exec(context.Background(), `
		DELETE FROM
			departments
		WHERE
			id = $1
	`, id)
	if err != nil {
		errmsg("DeleteDepartment Exec", err)
	}
	return err
}

func departmentCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			departments (
				id bigserial primary key,
				name text,
				note text,
				created_at TIMESTAMP without time zone,
				updated_at TIMESTAMP without time zone default now(),
				UNIQUE(name)
			)
	`
	_, err := pool.Exec(context.Background(), str)
	if err != nil {
		errmsg("departmentCreateTable exec", err)
	}
	return err
}
