package edc

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
		return department, nil
	}
	err := e.db.Model(&department).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetDepartment select", err)
	}
	return department, err
}

// GetDepartmentList - get all department for list
func (e *Edb) GetDepartmentList() ([]Department, error) {
	var departments []Department
	err := e.db.Model(&departments).
		Column("id", "name", "note").
		Order("name ASC").
		Select()
	if err != nil {
		errmsg("GetDepartmentList select", err)
	}
	return departments, err
}

// GetDepartmentSelect - get all department for select
func (e *Edb) GetDepartmentSelect() ([]SelectItem, error) {
	var departments []SelectItem
	err := e.db.Model(&Department{}).
		Column("id", "name").
		Order("name ASC").
		Select(&departments)
	if err != nil {
		errmsg("GetDepartmentSelect select", err)
	}
	return departments, err
}

// CreateDepartment - create new department
func (e *Edb) CreateDepartment(department Department) (int64, error) {
	err := e.db.Insert(&department)
	if err != nil {
		errmsg("CreateDepartment insert", err)
	}
	return department.ID, nil
}

// UpdateDepartment - save department changes
func (e *Edb) UpdateDepartment(department Department) error {
	err := e.db.Update(&department)
	if err != nil {
		errmsg("UpdateDepartment update", err)
	}
	return err
}

// DeleteDepartment - delete department by id
func (e *Edb) DeleteDepartment(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Model(&Department{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteDepartment delete", err)
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
				updated_at TIMESTAMP without time zone default now(),
				UNIQUE(name)
			)
	`
	_, err := e.db.Exec(str)
	if err != nil {
		errmsg("departmentCreateTable exec", err)
	}
	return err
}
