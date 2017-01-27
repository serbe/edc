package edc

import (
	"fmt"
	"log"
)

// Practice - struct for practice
type Practice struct {
	ID             int64   `sql:"id" json:"id"`
	Company        Company `sql:"-"`
	CompanyID      int64   `sql:"company_id, null" json:"company_id"`
	Kind           Kind    `sql:"-"`
	KindID         int64   `sql:"kind_id, null" json:"kind_id"`
	Topic          string  `sql:"topic, null" json:"topic"`
	DateOfPractice string  `sql:"date_of_practice, null" json:"date_of_practice"`
	DateStr        string  `sql:"-" json:"date_str"`
	Note           string  `sql:"note, null" json:"note"`
	CreatedAt      string  `sql:"created_at" json:"created_at"`
	UpdatedAt      string  `sql:"updated_at" json:"updated_at"`
}

// GetPractice - get one practice by id
func (e *Edb) GetPractice(id int64) (Practice, error) {
	if id == 0 {
		return Practice{}, nil
	}
	stmt, err := e.db.Prepare(`SELECT
		id,
		company_id,
		kind_id,
		topic,
		date_of_practice,
		note
	FROM
		practices
	WHERE id = $1`)
	if err != nil {
		log.Println("GetPractice e.db.Prepare ", err)
		return Practice{}, err
	}
	row := stmt.QueryRow(id)
	practice, err := scanPractice(row)
	return practice, err
}

// GetPracticeList - get all practices for list
func (e *Edb) GetPracticeList() ([]Practice, error) {
	rows, err := e.db.Query(`SELECT
		p.id,
		p.company_id,
		c.name AS company_name,
		k.name AS kind_name,
		p.topic,
		p.date_of_practice
	FROM
		practices AS p
	LEFT JOIN
		companies AS c ON c.id = p.company_id
	LEFT JOIN
		kinds AS k ON k.id = p.kind_id
	ORDER BY
		date_of_practice DESC`)
	if err != nil {
		log.Println("GetPracticeAll e.db.Query ", err)
		return []Practice{}, err
	}
	practices, err := scanPractices(rows, "list")
	return practices, err
}

// GetPracticeCompany - get all practices of company
func (e *Edb) GetPracticeCompany(id int64) ([]Practice, error) {
	if id == 0 {
		return []Practice{}, nil
	}
	stmt, err := e.db.Prepare(`SELECT
		p.id,
		k.name AS kind_name,
		p.topic,
		p.date_of_practice
	FROM
		practices AS p
	LEFT JOIN
		kinds AS k ON k.id = p.kind_id
	WHERE
	    p.company_id = $1
	ORDER BY
		date_of_practice`)
	if err != nil {
		log.Println("GetPracticeCompany e.db.Prepare ", err)
		return []Practice{}, err
	}
	rows, err := stmt.Query(id)
	if err != nil {
		log.Println("GetPracticeCompany stmt.Query ", err)
		return []Practice{}, err
	}
	practices, err := scanPractices(rows, "company")
	return practices, err
}

// GetPracticeNear - get 10 nearest practices
func (e *Edb) GetPracticeNear() ([]Practice, error) {
	rows, err := e.db.Query(`SELECT
		p.id,
		c.name AS company_name,
		k.name AS kind_name,
		p.topic,
		p.date_of_practice
	FROM
		practices AS p
	LEFT JOIN
		companies AS c ON c.id = p.company_id
	LEFT JOIN
		kinds AS k ON k.id = p.kind_id
	WHERE
		p.date_of_practice > now()
	ORDER BY
		date_of_practice
	LIMIT 10`)
	if err != nil {
		log.Println("GetPracticeNear e.db.Query ", err)
		return []Practice{}, err
	}
	practices, err := scanPractices(rows, "near")
	return practices, err
}

// CreatePractice - create new practice
func (e *Edb) CreatePractice(practice Practice) (int64, error) {
	stmt, err := e.db.Prepare(`
		INSERT INTO
			practices (
				company_id,
				kind_id,
				topic,
				date_of_practice,
				note,
				created_at
			) VALUES (
				$1,
				$2,
				$3,
				$4,
				$5,
				now()
			)
		RETURNING id
	`)
	if err != nil {
		log.Println("CreatePractice e.db.Prepare ", err)
		return 0, err
	}
	err = stmt.QueryRow(i2n(practice.CompanyID), i2n(practice.KindID), s2n(practice.Topic), sd2n(practice.DateOfPractice), s2n(practice.Note)).Scan(&practice.ID)
	return practice.ID, err
}

// UpdatePractice - save practice changes
func (e *Edb) UpdatePractice(practice Practice) error {
	stmt, err := e.db.Prepare(`
		UPDATE
			practices
		SET
			company_id = $2,
			kind_id = $3,
			topic = $4,
			date_of_practice = $5,
			note = $6,
			updated_at = now()
		WHERE
			id = $1
	`)
	if err != nil {
		log.Println("UpdatePractice e.db.Prepare ", err)
		return err
	}
	_, err = stmt.Exec(practice.ID, i2n(practice.CompanyID), i2n(practice.KindID), s2n(practice.Topic), sd2n(practice.DateOfPractice), s2n(practice.Note))
	return err
}

// DeletePractice - delete practice by id
func (e *Edb) DeletePractice(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec(`
		DELETE FROM
			practices
		WHERE
			id = $1
	`, id)
	if err != nil {
		log.Println("DeletePractice e.db.Exec: ", id, err)
		return fmt.Errorf("DeletePractice e.db.Exec: %s", err)
	}
	return err
}

func (e *Edb) practiceCreateTable() error {
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
				updated_at TIMESTAMP without time zone
			)
	`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("practiceCreateTable e.db.Exec: ", err)
		return fmt.Errorf("practiceCreateTable e.db.Exec: %s", err)
	}
	return err
}
