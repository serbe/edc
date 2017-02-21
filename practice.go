package edc

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

// PracticeList is struct for practice list
type PracticeList struct {
	ID             int64  `sql:"id" json:"id"`
	CompanyID      int64  `sql:"company_id, null" json:"company_id"`
	CompanyName    string `sql:"company_name, null" json:"company_name"`
	KindName       string `sql:"kind_name, null" json:"kind_name"`
	Topic          string `sql:"topic, null" json:"topic"`
	DateOfPractice string `sql:"date_of_practice, null" json:"date_of_practice"`
	DateStr        string `sql:"-" json:"date_str"`
}

// GetPractice - get one practice by id
func (e *Edb) GetPractice(id int64) (Practice, error) {
	var practice Practice
	if id == 0 {
		return practice, nil
	}
	err := e.db.Model(&practice).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetPractice select", err)
		return practice, err
	}
	return practice, err
}

// GetPracticeList - get all practices for list
func (e *Edb) GetPracticeList() ([]PracticeList, error) {
	var practices []PracticeList
	_, err := e.db.Query(&practices, `
	SELECT
		p.id,
		p.company_id,
		c.name AS company_name,
		k.name AS kind_name,
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
		errmsg("GetPracticeList query", err)
	}
	for i := range practices {
		practices[i].DateStr = setStrMonth(practices[i].DateOfPractice)
	}
	return practices, err
}

// GetPracticeCompany - get all practices of company
func (e *Edb) GetPracticeCompany(id int64) ([]Practice, error) {
	var practices []Practice
	if id == 0 {
		return practices, nil
	}
	err := e.db.Model(&practices).
		Where("company_id = ?", id).
		Select()
	for i := range practices {
		practices[i].DateStr = setStrMonth(practices[i].DateOfPractice)
	}
	if err != nil {
		errmsg("GetPracticeCompany select", err)
	}
	return practices, err
}

// GetPracticeNear - get 10 nearest practices
func (e *Edb) GetPracticeNear() ([]PracticeList, error) {
	var practices []PracticeList
	_, err := e.db.Query(&practices, `
	SELECT
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
	WHERE
		p.date_of_practice > now()
	ORDER BY
		date_of_practice
	LIMIT 10`)
	for i := range practices {
		practices[i].DateStr = setStrMonth(practices[i].DateOfPractice)
	}
	if err != nil {
		errmsg("GetPracticeNear query", err)
	}
	return practices, err
}

// CreatePractice - create new practice
func (e *Edb) CreatePractice(practice Practice) (int64, error) {
	err := e.db.Insert(&practice)
	if err != nil {
		errmsg("CreatePractice insert", err)
	}
	return practice.ID, err
}

// UpdatePractice - save practice changes
func (e *Edb) UpdatePractice(practice Practice) error {
	err := e.db.Update(&practice)
	if err != nil {
		errmsg("UpdatePractice update", err)
	}
	return err
}

// DeletePractice - delete practice by id
func (e *Edb) DeletePractice(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Model(&Practice{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeletePractice delete", err)
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
				updated_at TIMESTAMP without time zone default now()
			)
	`
	_, err := e.db.Exec(str)
	if err != nil {
		errmsg("practiceCreateTable exec", err)
	}
	return err
}
