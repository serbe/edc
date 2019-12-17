package edc

import "context"

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
	err := pool.QueryRow(context.Background(), &department).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetDepartment select", err)
	}
	return department, err
}

// DepartmentListGet - get all department for list
func DepartmentListGet() ([]DepartmentList, error) {
	var departments []DepartmentList
	err := pool.QueryRow(context.Background(), &Department{}).
		Column("id", "name", "note").
		Order("name ASC").
		Select(&departments)
	if err != nil {
		errmsg("GetDepartmentList select", err)
	}
	return departments, err
}

// DepartmentSelectGet - get all department for select
func DepartmentSelectGet() ([]SelectItem, error) {
	var departments []SelectItem
	err := pool.QueryRow(context.Background(), &Department{}).
		Column("id", "name").
		Order("name ASC").
		Select(&departments)
	if err != nil {
		errmsg("GetDepartmentSelectAll select", err)
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
	return departments, err
}

// DepartmentInsert - create new department
func DepartmentInsert(department Department) (int64, error) {
	err := pool.Insert(&department)
	if err != nil {
		errmsg("CreateDepartment insert", err)
	}
	return department.ID, nil
}

// DepartmentUpdate - save department changes
func DepartmentUpdate(department Department) error {
	err := pool.Update(&department)
	if err != nil {
		errmsg("UpdateDepartment update", err)
	}
	return err
}

// DepartmentDelete - delete department by id
func DepartmentDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.QueryRow(context.Background(), &Department{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteDepartment delete", err)
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
				updated_at
 TIMESTAMP without time zone default now(),
				UNIQUE(name)
			)
	`
	_, err := pool.Exec(str)
	if err != nil {
		errmsg("departmentCreateTable exec", err)
	}
	return err
}
