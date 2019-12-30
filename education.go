package edc

import (
	"context"
	"time"
)

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
	education.ID = id
	err := pool.QueryRow(context.Background(), `
		SELECT
			contact_id,
			start_date,
			end_date,
			post_id,
			note,
			created_at,
			updated_at
		FROM
			educations
		WHERE
			id = $1
	`, id).Scan(&education.ContactID, &education.StartDate, &education.EndDate, &education.PostID, &education.Note, &education.CreatedAt, &education.UpdatedAt)
	if err != nil {
		errmsg("EducationGet QueryRow", err)
	}
	return education, err
}

// EducationListGet - get all education for list
func EducationListGet() ([]EducationList, error) {
	var educations []EducationList
	rows, err := pool.Query(context.Background(), `
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
	if err != nil {
		errmsg("EducationListGet Query", err)
		return educations, err
	}
	for rows.Next() {
		var education EducationList
		err := rows.Scan(&education.ID, &education.ContactID, &education.ContactName, &education.StartDate,
			&education.EndDate, &education.PostID, &education.PostName, &education.Note)
		if err != nil {
			errmsg("EducationListGet Scan", err)
			return educations, err
		}
		educations = append(educations, education)
	}
	// for i := range educations {
	// 	educations[i].StartStr = setStrMonth(educations[i].StartDate)
	// 	educations[i].EndStr = setStrMonth(educations[i].EndDate)
	// }
	return educations, rows.Err()
}

// EducationNearGet - get 10 nearest educations
func EducationNearGet() ([]EducationShort, error) {
	var educations []EducationShort
	rows, err := pool.Query(context.Background(), `
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
		errmsg("EducationNearGet Query", err)
	}
	for rows.Next() {
		var education EducationShort
		err := rows.Scan(&education.ID, &education.ContactID, &education.ContactName, &education.StartDate)
		if err != nil {
			errmsg("EducationNearGet Scan", err)
			return educations, err
		}
		educations = append(educations, education)
	}
	return educations, rows.Err()
}

// EducationInsert - create new education
func EducationInsert(education Education) (int64, error) {
	err := pool.QueryRow(context.Background(), `
		INSERT INTO educations
		(
			contact_id,
			start_date,
			end_date,
			post_id,
			note,
			created_at,
			updated_at
		)
		VALUES
		(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7
		)
		RETURNING
			id
	`, education.ContactID, education.StartDate, education.EndDate, education.PostID,
		education.Note, time.Now(), time.Now()).Scan(&education.ID)
	if err != nil {
		errmsg("EducationInsert QueryRow", err)
	}
	return education.ID, err
}

// EducationUpdate - save changes to education
func EducationUpdate(education Education) error {
	_, err := pool.Exec(context.Background(), `
		UPDATE educations SET
			contact_id = $2,
			start_date = $3,
			end_date = $4,
			post_id = $5,
			note = $6,
			updated_at = $7
		WHERE
			id = $1
	`, education.ID, education.ContactID, education.StartDate, education.EndDate, education.PostID, education.Note, time.Now())
	if err != nil {
		errmsg("EducationUpdate update", err)
	}
	return err
}

// EducationDelete - delete education by id
func EducationDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.Exec(context.Background(), `
		DELETE FROM
			educations
		WHERE
			id = $1
	`, id)
	if err != nil {
		errmsg("EducationDelete Exec", err)
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
				updated_at TIMESTAMP without time zone default now()
			)
	`
	_, err := pool.Exec(context.Background(), str)
	if err != nil {
		errmsg("educationCreateTable exec", err)
	}
	return err
}
