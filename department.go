package edc

import "log"

// Department - struct for department
type Department struct {
	ID        int64  `sql:"id" json:"id"`
	Name      string `sql:"name" json:"name"`
	Note      string `sql:"note, null" json:"note"`
	CreatedAt string `sql:"created_at" json:"created_at"`
	UpdatedAt string `sql:"updated_at" json:"updated_at"`
}

// GetDepartment - get one department by id
func (e *Edb) GetDepartment(id int64) (Department, error) {
	var department Department
	if id == 0 {
		return Department{}, nil
	}

	err := e.db.Model(&department).Where(`id = ?`, id).Select()
	return department, err
}

// GetDepartmentList - get all department for list
func (e *Edb) GetDepartmentList() ([]Department, error) {
	var departments []Department
	_, err := e.db.Query(&departments, `
		SELECT
			id,
			name,
			note
		FROM
			departments
		ORDER BY
			name
		ASC
	`)
	if err != nil {
		log.Println("GetDepartmentList e.db.Query ", err)
		return []Department{}, err
	}
	return departments, err
}

// GetDepartmentSelect - get all department for select
func (e *Edb) GetDepartmentSelect() ([]SelectItem, error) {
	var departments []SelectItem
	rows, err := e.db.Query(&departments, `
		SELECT
			id,
			name
		FROM
			departments
		ORDER BY
			name ASC
	`)
	if err != nil {
		log.Println("GetDepartmentSelect e.db.Query ", err)
		return []SelectItem{}, err
	}
	return departments, err
}

// CreateDepartment - create new department
func (e *Edb) CreateDepartment(department Department) (int64, error) {
	err := e.db.Insert(&department)
	if err != nil {
		log.Println("CreateDepartment e.db.Insert ", err)
		return 0, err
	}
	return department.ID, nil
}

// UpdateDepartment - save department changes
func (e *Edb) UpdateDepartment(department Department) error {
	err := e.db.Update(&department)
	if err != nil {
		log.Println("UpdateDepartment e.db.Update ", err)
		return err
	}
	return err
}

// DeleteDepartment - delete department by id
func (e *Edb) DeleteDepartment(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec(`
		DELETE FROM
			departments
		WHERE
			id = $1
	`, id)
	if err != nil {
		log.Println("DeleteDepartment e.db.Exec ", id, err)
	}
	return err
}

func (e *Edb) departmentCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			departments (
				id bigserial primary key,
				name text,
				note text,
				created_at TIMESTAMP without time zone,
				updated_at TIMESTAMP without time zone,
				UNIQUE(name)
			)
	`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("departmentCreateTable e.db.Exec ", err)
	}
	return err
}
