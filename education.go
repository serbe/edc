package edc

import "context"

// Education - struct for education
type Education struct {
	ID        int64  `sql:"id"         json:"id"         form:"id"         query:"id" `
	ContactID int64  `sql:"contact_id" json:"contact_id" form:"contact_id" query:"contact_id"`
	StartDate string `sql:"start_date" json:"start_date" form:"start_date" query:"start_date"`
	EndDate   string `sql:"end_date"   json:"end_date"   form:"end_date"   query:"end_date"`
	PostID    int64  `sql:"post_id"    json:"post_id"    form:"post_id"    query:"post_id"`
	Note      string `sql:"note"       json:"note"       form:"note"       query:"note"`
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
	Note        string `sql:"note"         json:"note"         form:"note"         query:"note"`
}

// EducationShort - short struct for education
type EducationShort struct {
	ID          int64  `sql:"id"           json:"id"           form:"id"           query:"id"`
	ContactID   int64  `sql:"contact_id"   json:"contact_id"   form:"contact_id"   query:"contact_id"`
	ContactName string `sql:"contact_name" json:"contact_name" form:"contact_name" query:"contact_name"`
	StartDate   string `sql:"start_date"   json:"start_date"   form:"start_date"   query:"start_date"`
}

// EducationGet - get education by id
func EducationGet(id int64) (Education, error) {
	var education Education
	if id == 0 {
		return education, nil
	}
	err := pool.QueryRow(context.Background(), &education).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetEducation select", err)
	}
	return education, err
}

// EducationListGet - get all education for list
func EducationListGet() ([]EducationList, error) {
	var educations []EducationList
	_, err := pool.Query(context.Background(), `
		SELECT
			e.id,
			e.contact_id,
			c.name AS contact_name,
			e.start_date,
			e.end_date,
			e.post_id,
			p.name AS post_name,
			e.note
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

// EducationNearGet - get 10 nearest educations
func EducationNearGet() ([]EducationShort, error) {
	var educations []EducationShort
	_, err := pool.Query(context.Background(), `
		SELECT
			e.id,
			e.contact_id,
			c.name AS contact_name,
			e.start_date
		FROM
			educations AS e
		LEFT JOIN
			contacts AS c ON c.id = e.contact_id
		WHERE
			e.start_date > TIMESTAMP 'now'::timestamp - '1 month'::interval
		ORDER BY
			start_date ASC
		LIMIT 10
	`)
	if err != nil {
		errmsg("GetEducationNear query", err)
	}
	return educations, err
}

// EducationSelectGet - get all education for select
func EducationSelectGet() ([]Education, error) {
	var educations []Education
	err := pool.QueryRow(context.Background(), &educations).
		C("id", "start_date", "end_date").
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
	return educations, err
}

// EducationInsert - create new education
func EducationInsert(education Education) (int64, error) {
	err := pool.Insert(&education)
	if err != nil {
		errmsg("CreateEducation insert", err)
	}
	return education.ID, err
}

// EducationUpdate - save changes to education
func EducationUpdate(education Education) error {
	err := pool.Update(&education)
	if err != nil {
		errmsg("UpdateEducation update", err)
	}
	return err
}

// EducationDelete - delete education by id
func EducationDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.QueryRow(context.Background(), &Education{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteEducation delete", err)
	}
	return err
}

func educationCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			educations (
				id bigserial primary key,
				start_date date,
				end_date date,
				note text,
				post_id bigint,
				created_at TIMESTAMP without time zone,
				updated_at
 TIMESTAMP without time zone default now()
			)
	`
	_, err := pool.Exec(str)
	if err != nil {
		errmsg("educationCreateTable exec", err)
	}
	return err
}
