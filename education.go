package edc

// Education - struct for education
type Education struct {
	ID        int64  `sql:"id"         json:"id"         form:"id"         query:"id" `
	ContactID int64  `sql:"contact_id" json:"contact_id" form:"contact_id" query:"contact_id"`
	StartDate string `sql:"start_date" json:"start_date" form:"start_date" query:"start_date"`
	EndDate   string `sql:"end_date"   json:"end_date"   form:"end_date"   query:"end_date"`
	StartStr  string `sql:"-"          json:"start_str"  form:"start_str"  query:"start_str"`
	EndStr    string `sql:"-"          json:"end_str"    form:"end_str"    query:"end_str"`
	PostID    int64  `sql:"post_id"    json:"post_id"    form:"post_id"    query:"post_id"`
	Note      string `sql:"note,null"  json:"note"       form:"note"       query:"note"`
	CreatedAt string `sql:"created_at" json:"-"`
	UpdatedAt string `sql:"updated_at" json:"-"`
}

// EducationList - struct for list of education
type EducationList struct {
	ID          int64  `sql:"id"           json:"id"           form:"id"           query:"id"`
	ContactID   int64  `sql:"contact_id"   json:"contact_id"   form:"contact_id"   query:"contact_id"`
	ContactName string `sql:"contact_name" json:"contact_name" form:"contact_name" query:"contact_name"`
	StartDate   string `sql:"start_date"   json:"start_date"   form:"start_date"   query:"start_date"`
	EndDate     string `sql:"end_date"     json:"end_date"     form:"end_date"     query:"end_date"`
	StartStr    string `sql:"-"            json:"start_str"    form:"start_str"    query:"start_str"`
	EndStr      string `sql:"-"            json:"end_str"      form:"end_str"      query:"end_str"`
	PostID      int64  `sql:"post_id"      json:"post_id"      form:"post_id"      query:"post_id"`
	PostName    string `sql:"post_name"    json:"post_name"    form:"post_name"    query:"post_name"`
	Note        string `sql:"note,null"    json:"note"         form:"note"         query:"note"`
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

// GetEducationListAll - get all education for list
func (e *Edb) GetEducationListAll() ([]EducationList, error) {
	var educations []EducationList
	_, err := e.db.Query(&educations, `
		SELECT
			e.id,
			e.start_date,
			e.end_date,
			e.contact_id,
			c.name AS contact_name,
			e.post_id,
			p.name AS post_name
		FROM
			educations AS e
		LEFT JOIN
			contacts AS c ON c.id = e.contact_id
		LEFT JOIN
			posts AS p ON p.id = e.post_id
		ORDER BY
			start_date DESC
	`)
	for i := range educations {
		educations[i].StartStr = setStrMonth(educations[i].StartDate)
		educations[i].EndStr = setStrMonth(educations[i].EndDate)
	}
	return educations, err
}

// GetEducationNear - get 10 nearest educations
func (e *Edb) GetEducationNear() ([]EducationList, error) {
	var educations []EducationList
	_, err := e.db.Query(&educations, `
		SELECT
			e.id,
			e.start_date,
			e.end_date,
			e.contact_id,
			c.name AS contact_name
		FROM
			educations AS e
		LEFT JOIN
			contacts AS c ON c.id = e.contact_id
		WHERE
			e.start_date > TIMESTAMP 'now'::timestamp - '1 month'::interval
		ORDER BY
			start_date DESC
		LIMIT 10
	`)
	for i := range educations {
		educations[i].StartStr = setStrMonth(educations[i].StartDate)
	}
	if err != nil {
		errmsg("GetEducationNear query", err)
	}
	return educations, err
}

// // GetEducationSelectAll - get all education for select
// func (e *Edb) GetEducationSelectAll() ([]Education, error) {
// 	var educations []Education
// 	err := e.db.Model(&educations).
// 		C("id", "start_date", "end_date").
// 		Order("start_date").
// 		Select()
// 	if err != nil {
// 		errmsg("GetEducationSelectAll select", err)
// 		return educations, err
// 	}
// 	for i := range educations {
// 		educations[i].StartStr = setStrMonth(educations[i].StartDate)
// 		educations[i].EndStr = setStrMonth(educations[i].EndDate)
// 	}
// 	return educations, err
// }

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
				post_id bigint,
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
