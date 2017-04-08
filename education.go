package edc

// Education - struct for education
type Education struct {
	ID        int64  `sql:"id"         json:"id" `
	StartDate string `sql:"start_date" json:"start_date"`
	EndDate   string `sql:"end_date"   json:"end_date"`
	StartStr  string `sql:"-"          json:"start_str"`
	EndStr    string `sql:"-"          json:"end_str"`
	Note      string `sql:"note, null" json:"note"`
	CreatedAt string `sql:"created_at" json:"-"`
	UpdatedAt string `sql:"updated_at" json:"-"`
}

// GetEducation - get education by id
func (e *Edb) GetEducation(id int64) (Education, error) {
	var education Education
	if id == 0 {
		return education, nil
	}
	err := e.db.Model(&education).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetEducation select", err)
	}
	return education, err
}

// GetEducationList - get all education for list
func (e *Edb) GetEducationList() ([]Education, error) {
	var educations []Education
	err := e.db.Model(&educations).
		Column("id", "start_date", "end_date", "note").
		Order("start_date").
		Select()
	if err != nil {
		errmsg("GetEducationList select", err)
		return educations, err
	}
	for i := range educations {
		educations[i].StartStr = setStrMonth(educations[i].StartDate)
		educations[i].EndStr = setStrMonth(educations[i].EndDate)
	}
	return educations, err
}

// GetEducationSelectAll - get all education for select
func (e *Edb) GetEducationSelectAll() ([]Education, error) {
	var educations []Education
	err := e.db.Model(&educations).
		Where("id", "start_date", "end_date").
		Order("start_date").
		Select()
	if err != nil {
		errmsg("GetEducationSelectAll select", err)
		return educations, err
	}
	for i := range educations {
		educations[i].StartStr = setStrMonth(educations[i].StartDate)
		educations[i].EndStr = setStrMonth(educations[i].EndDate)
	}
	return educations, err
}

// CreateEducation - create new education
func (e *Edb) CreateEducation(education Education) (int64, error) {
	err := e.db.Insert(&education)
	if err != nil {
		errmsg("CreateEducation insert", err)
	}
	return education.ID, err
}

// UpdateEducation - save changes to education
func (e *Edb) UpdateEducation(education Education) error {
	err := e.db.Update(&education)
	if err != nil {
		errmsg("UpdateEducation update", err)
	}
	return err
}

// DeleteEducation - delete education by id
func (e *Edb) DeleteEducation(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Model(&Education{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteEducation delete", err)
	}
	return err
}

func (e *Edb) educationCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			educations (
				id bigserial primary key,
				start_date date,
				end_date date,
				note text,
				created_at TIMESTAMP without time zone,
				updated_at TIMESTAMP without time zone default now()
			)
	`
	_, err := e.db.Exec(str)
	if err != nil {
		errmsg("educationCreateTable exec", err)
	}
	return err
}
