package edc

import "log"

// Education - struct for education
type Education struct {
	ID        int64  `sql:"id" json:"id" `
	StartDate string `sql:"start_date" json:"start_date"`
	EndDate   string `sql:"end_date" json:"end_date"`
	StartStr  string `sql:"-" json:"start_str"`
	EndStr    string `sql:"-" json:"end_str"`
	Note      string `sql:"note, null" json:"note"`
	CreatedAt string `sql:"created_at" json:"created_at"`
	UpdatedAt string `sql:"updated_at" json:"updated_at"`
}

// GetEducation - get education by id
func (e *Edb) GetEducation(id int64) (Education, error) {
	var education Education
	if id == 0 {
		return education, nil
	}
	err := e.db.Model(&education).Where(`id = ?`, id).Select()
	if err != nil {
		errmsg("GetEducation select", err)
	}
	return education, err
}

// GetEducationList - get all education for list
func (e *Edb) GetEducationList() ([]Education, error) {
	var educations []Education
	err := e.db.Model(&educations).
		Where("educations.id", "educations.start_date", "educations.end_date", "educations.note").
		Order("educations.start_date").
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

// GetEducationSelect - get all education for select
func (e *Edb) GetEducationSelect() ([]Education, error) {
	var educations []Education
	err := e.db.Model(&educations).
		Where("educations.id", "educations.start_date", "educations.end_date").
		Order("educations.start_date").
		Select()
	if err != nil {
		errmsg("GetEducationSelect select", err)
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
		log.Println("CreateEducation e.db.Insert ", err)
		return 0, err
	}
	return education.ID, err
}

// UpdateEducation - save changes to education
func (e *Edb) UpdateEducation(education Education) error {
	stmt, err := e.db.Prepare(`
		UPDATE
			educations
		SET
			start_date = $2,
			end_date = $3,
			note = $4,
			updated_at = now()
		WHERE
			id = $1
	`)
	if err != nil {
		log.Println("UpdateEducation e.db.Prepare ", err)
		return err
	}
	_, err = stmt.Exec(education.ID, sd2n(education.StartDate), sd2n(education.EndDate), s2n(education.Note))
	return err
}

// DeleteEducation - delete education by id
func (e *Edb) DeleteEducation(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec(`
		DELETE FROM
			educations
		WHERE
			id = $1
	`, id)
	if err != nil {
		log.Println("DeleteEducation ", id, err)
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
				updated_at TIMESTAMP without time zone
			)
	`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("educationCreateTable ", err)
	}
	return err
}
