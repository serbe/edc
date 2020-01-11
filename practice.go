package edc

import (
	"context"
	"time"
)

// Practice - struct for practice
type Practice struct {
	ID             int64  `sql:"id"               json:"id"               form:"id"               query:"id"`
	CompanyID      int64  `sql:"company_id"       json:"company_id"       form:"company_id"       query:"company_id"`
	KindID         int64  `sql:"kind_id"          json:"kind_id"          form:"kind_id"          query:"kind_id"`
	Topic          string `sql:"topic"            json:"topic"            form:"topic"            query:"topic"`
	DateOfPractice string `sql:"date_of_practice" json:"date_of_practice" form:"date_of_practice" query:"date_of_practice"`
	Note           string `sql:"note"             json:"note"             form:"note"             query:"note"`
	CreatedAt      string `sql:"created_at"       json:"-"`
	UpdatedAt      string `sql:"updated_at"       json:"-"`
}

// PracticeList is struct for practice list
type PracticeList struct {
	ID             int64  `sql:"id"               json:"id"               form:"id"               query:"id"`
	CompanyID      int64  `sql:"company_id"       json:"company_id"       form:"company_id"       query:"company_id"`
	CompanyName    string `sql:"company_name"     json:"company_name"     form:"company_name"     query:"company_name"`
	KindID         int64  `sql:"kind_id"          json:"kind_id"          form:"kind_id"          query:"kind_id"`
	KindName       string `sql:"-"                json:"kind_name"        form:"kind_name"        query:"kind_name"`
	KindShortName  string `sql:"-"                json:"kind_short_name"  form:"kind_short_name"  query:"kind_short_name"`
	Topic          string `sql:"topic"            json:"topic"            form:"topic"            query:"topic"`
	DateOfPractice string `sql:"date_of_practice" json:"date_of_practice" form:"date_of_practice" query:"date_of_practice"`
	DateStr        string `sql:"-"                json:"date_str"         form:"date_str"         query:"date_str"`
}

// PracticeShort - short struct for practice
type PracticeShort struct {
	ID             int64  `sql:"id"               json:"id"               form:"id"               query:"id"`
	CompanyID      int64  `sql:"company_id"       json:"company_id"       form:"company_id"       query:"company_id"`
	CompanyName    string `sql:"company_name"     json:"company_name"     form:"company_name"     query:"company_name"`
	KindID         int64  `sql:"kind_id"          json:"kind_id"          form:"kind_id"          query:"kind_id"`
	KindShortName  string `sql:"-"                json:"kind_short_name"  form:"kind_short_name"  query:"kind_short_name"`
	DateOfPractice string `sql:"date_of_practice" json:"date_of_practice" form:"date_of_practice" query:"date_of_practice"`
}

// PracticeGet - get one practice by id
func PracticeGet(id int64) (Practice, error) {
	var practice Practice
	if id == 0 {
		return practice, nil
	}
	practice.ID = id
	err := pool.QueryRow(context.Background(), `
		SELECT
			company_id,
			kind_id,
			topic,
			date_of_practice,
			note,
			created_at,
			updated_at
		FROM
			practices
		WHERE
			id = $1
	`, id).Scan(&practice.CompanyID, &practice.KindID, &practice.Topic, &practice.DateOfPractice, &practice.Note,
		&practice.CreatedAt, &practice.UpdatedAt)
	if err != nil {
		errmsg("PracticeGet QueryRow", err)
		return practice, err
	}
	return practice, err
}

// PracticeListGet - get all practices for list
func PracticeListGet() ([]PracticeList, error) {
	var practices []PracticeList
	_, err := pool.Query(context.Background(), `
		SELECT
			p.id,
			p.company_id,
			c.name AS company_name,
			k.name AS kind_name,
			k.short_name AS kind_short_name,
			p.date_of_practice,
			p.topic
		FROM
			practices AS p
		LEFT JOIN
			companies AS c ON c.id = p.company_id
		LEFT JOIN
			kinds AS k ON k.id = p.kind_id
		ORDER BY
			date_of_practice DESC`)
	if err != nil {
		errmsg("GetPracticeList query", err)
	}
	for i := range practices {
		practices[i].DateStr = setStrMonth(practices[i].DateOfPractice)
	}
	return practices, err
}

// PracticeCompanyGet - get all practices of company
func PracticeCompanyGet(id int64) ([]PracticeList, error) {
	var practices []PracticeList
	if id == 0 {
		return practices, nil
	}
	rows, err := pool.Query(context.Background(), `
		SELECT
			p.id,
			p.company_id,
			c.name AS company_name,
			p.kind_id,
			k.name AS kind_name,
			k.short_name AS kind_short_name,
			p.date_of_practice,
			p.topic
		FROM
			practices AS p
		LEFT JOIN
			companies AS c ON c.id = p.company_id
		LEFT JOIN
			kinds AS k ON k.id = p.kind_id
		WHERE
			p.company_id = $1
		ORDER BY
			date_of_practice DESC
	`, id)
	if err != nil {
		errmsg("GetPracticeCompany query", err)
		return practices, err
	}
	for rows.Next() {
		var practice PracticeList
		err := rows.Scan(&practice.ID, &practice.CompanyID, &practice.CompanyName,
			&practice.KindID, &practice.KindName, &practice.KindShortName, &practice.DateOfPractice, &practice.Topic)
		if err != nil {
			errmsg("GetPracticeCompany select", err)
			return practices, err
		}
		practice.DateStr = setStrMonth(practice.DateOfPractice)
		practices = append(practices, practice)
	}
	return practices, rows.Err()
}

// PracticeNearGet - get 10 nearest practices
func PracticeNearGet() ([]PracticeShort, error) {
	var practices []PracticeShort
	_, err := pool.Query(context.Background(), `
		SELECT
			p.id,
			p.company_id,
			c.name AS company_name,
			p.kind_id,
			k.short_name AS kind_short_name,
			p.date_of_practice
		FROM
			practices AS p
		LEFT JOIN
			companies AS c ON c.id = p.company_id
		LEFT JOIN
			kinds AS k ON k.id = p.kind_id
		WHERE
			p.date_of_practice > TIMESTAMP 'now'::timestamp - '1 month'::interval
		ORDER BY
			date_of_practice ASC
		LIMIT 10`)
	if err != nil {
		errmsg("GetPracticeNear query", err)
	}
	return practices, err
}

// PracticeInsert - create new practice
func PracticeInsert(practice Practice) (int64, error) {
	err := pool.QueryRow(context.Background(), `
		INSERT INTO practices
		(
			company_id,
			kind_id,
			topic,
			date_of_practice,
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
	`, practice.CompanyID, practice.KindID, practice.Topic, practice.DateOfPractice,
		practice.Note, time.Now(), time.Now()).Scan(&practice.ID)
	if err != nil {
		errmsg("PracticeInsert QueryRow", err)
	}
	return practice.ID, err
}

// PracticeUpdate - save practice changes
func PracticeUpdate(practice Practice) error {
	_, err := pool.Exec(context.Background(), `
		UPDATE practices SET
			company_id = $2,
			kind_id = $3,
			topic = $4,
			date_of_practice = $5,
			note = $6,
			updated_at = $7
		WHERE
			id = $1
	`, practice.ID, practice.CompanyID, practice.KindID, practice.Topic, practice.DateOfPractice,
		practice.Note, time.Now())
	if err != nil {
		errmsg("PracticeUpdate Exec", err)
	}
	return err
}

// PracticeDelete - delete practice by id
func PracticeDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.Exec(context.Background(), `
		DELETE FROM
			practices
		WHERE
			id = $1
	`, id)
	if err != nil {
		errmsg("DeletePractice Exec", err)
	}
	return err
}

func practiceCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			practices (
				id bigserial primary key,
				company_id bigint,
				kind_id bigint,
				topic text,
				date_of_practice date,
				note text,
				created_at TIMESTAMP without time zone,
				updated_at TIMESTAMP without time zone default now()
			)
	`
	_, err := pool.Exec(context.Background(), str)
	if err != nil {
		errmsg("practiceCreateTable exec", err)
	}
	return err
}
