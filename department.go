package edc

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

// GetDepartment - get one department by id
func GetDepartment(id int64) (Department, error) {
	var department Department
	if id == 0 {
		return department, nil
	}
	err := pool.Model(&department).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetDepartment select", err)
	}
	return department, err
}

// GetDepartmentList - get department for list by id
func GetDepartmentList(id int64) (DepartmentList, error) {
	var department DepartmentList
	err := pool.Model(&Department{}).
		Column("id", "name", "note").
		Where("id = ?", id).
		Select(&department)
	if err != nil {
		errmsg("GetDepartmentList select", err)
	}
	return department, err
}

// GetDepartmentListAll - get all department for list
func GetDepartmentListAll() ([]DepartmentList, error) {
	var departments []DepartmentList
	err := pool.Model(&Department{}).
		Column("id", "name", "note").
		Order("name ASC").
		Select(&departments)
	if err != nil {
		errmsg("GetDepartmentList select", err)
	}
	return departments, err
}

// GetDepartmentSelect - get department for select
func GetDepartmentSelect(id int64) (SelectItem, error) {
	var department SelectItem
	if id == 0 {
		return department, nil
	}
	err := pool.Model(&Department{}).
		Column("id", "name").
		Where("id = ?", id).
		Select(&department)
	if err != nil {
		errmsg("GetDepartmentSelect select", err)
	}
	return department, err
}

// GetDepartmentSelectAll - get all department for select
func GetDepartmentSelectAll() ([]SelectItem, error) {
	var departments []SelectItem
	err := pool.Model(&Department{}).
		Column("id", "name").
		Order("name ASC").
		Select(&departments)
	if err != nil {
		errmsg("GetDepartmentSelectAll select", err)
	}
	return departments, err
}

// CreateDepartment - create new department
func CreateDepartment(department Department) (int64, error) {
	err := pool.Insert(&department)
	if err != nil {
		errmsg("CreateDepartment insert", err)
	}
	return department.ID, nil
}

// UpdateDepartment - save department changes
func UpdateDepartment(department Department) error {
	err := pool.Update(&department)
	if err != nil {
		errmsg("UpdateDepartment update", err)
	}
	return err
}

// DeleteDepartment - delete department by id
func DeleteDepartment(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.Model(&Department{}).
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
				updated_at TIMESTAMP without time zone default now(),
				UNIQUE(name)
			)
	`
	_, err := pool.Exec(str)
	if err != nil {
		errmsg("departmentCreateTable exec", err)
	}
	return err
}
