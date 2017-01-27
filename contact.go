package edc

import "log"

// Contact is struct for contact
type Contact struct {
	ID           int64       `sql:"id" json:"id"`
	Name         string      `sql:"name" json:"name"`
	Company      Company     `sql:"-"`
	CompanyID    int64       `sql:"company_id, null" json:"company_id"`
	Department   Department  `sql:"-"`
	DepartmentID int64       `sql:"department_id, null" json:"department_id"`
	Post         Post        `sql:"-"`
	PostID       int64       `sql:"post_id, null" json:"post_id"`
	PostGO       Post        `sql:"-"`
	PostGOID     int64       `sql:"post_go_id, null" json:"post_go_id"`
	Rank         Rank        `sql:"-"`
	RankID       int64       `sql:"rank_id, null" json:"rank_id"`
	Birthday     string      `sql:"birthday, null" json:"birthday"`
	Note         string      `sql:"note, null" json:"note"`
	Emails       []Email     `sql:"-"`
	Phones       []Phone     `sql:"-"`
	Faxes        []Phone     `sql:"-"`
	Educations   []Education `sql:"-"`
	CreatedAt    string      `sql:"created_at" json:"created_at"`
	UpdatedAt    string      `sql:"updated_at" json:"updated_at"`
}

// ContactList is struct for contact list
type ContactList struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	CompanyID   int64    `json:"company_id"`
	CompanyName string   `json:"company_name"`
	PostName    string   `json:"post_name"`
	Phones      []string `json:"phones"        pg:",array"`
	Faxes       []string `json:"faxes"         pg:",array"`
}

// ContactCompany is struct for company
type ContactCompany struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	DepartmentName string `json:"department_name"`
	PostName       string `json:"post_name"`
	PostGOName     string `json:"post_go_name"`
}

// GetContact - get one contact by id
func (e *Edb) GetContact(id int64) (Contact, error) {
	if id == 0 {
		return Contact{}, nil
	}
	var contact Contact
	// stmt, err := e.db.Query(`
	// 	SELECT
	// 		c.id,
	// 		c.name,
	// 		c.company_id,
	// 		c.department_id,
	// 		c.post_id,
	// 		c.post_go_id,
	// 		c.rank_id,
	// 		c.birthday,
	// 		c.note,
	// 		array_to_string(array_agg(DISTINCT e.email),',') AS email,
	// 		array_to_string(array_agg(DISTINCT p.phone),',') AS phone,
	// 		array_to_string(array_agg(DISTINCT f.phone),',') AS fax
	// 	FROM
	// 		contacts AS c
	// 	LEFT JOIN
	// 		emails AS e ON c.id = e.contact_id
	// 	LEFT JOIN
	// 		phones AS p ON c.id = p.contact_id AND p.fax = false
	// 	LEFT JOIN
	// 		phones AS f ON c.id = f.contact_id AND f.fax = true
	// 	WHERE
	// 		c.id = $1
	// 	GROUP BY
	// 		c.id
	// `)
	err := e.db.Model(&contact).Where("id = ?", id).Select()
	if err != nil {
		log.Println("GetContact e.db.Prepare ", err)
		return Contact{}, err
	}
	// contact.Educations = GetContactEducationscontacte.ID)
	return contact, err
}

// GetContactList - get all contacts for list
func (e *Edb) GetContactList() ([]ContactList, error) {
	var contacts []ContactList
	_, err := e.db.Query(&contacts, `
		SELECT
			c.id,
			c.name,
			co.id AS company_id,
			co.name AS company_name,
			po.name AS post_name,
			array_to_string(array_agg(DISTINCT ph.phone),',') AS phone,
			array_to_string(array_agg(DISTINCT f.phone),',') AS fax
		FROM
			contacts AS c
		LEFT JOIN
			companies AS co ON c.company_id = co.id
		LEFT JOIN
			posts AS po ON c.post_id = po.id
		LEFT JOIN
			phones AS ph ON c.id = ph.contact_id AND ph.fax = false
		LEFT JOIN
			phones AS f ON c.id = f.contact_id AND f.fax = true
		GROUP BY
			c.id,
			co.id,
			po.name
		ORDER BY
			name ASC
	`)
	if err != nil {
		log.Println("GetContactList e.db.Query ", err)
		return []ContactList{}, err
	}
	return contacts, err
}

// GetContactSelect - get all contacts for select
func (e *Edb) GetContactSelect() ([]SelectItem, error) {
	var contacts []SelectItem
	_, err := e.db.Query(&contacts, `
		SELECT
			c.id,
			c.name
		FROM
			contacts AS c
		ORDER BY
			name ASC
	`)
	if err != nil {
		log.Println("GetContactSelect e.db.Query ", err)
		return []SelectItem{}, err
	}
	return contacts, err
}

// GetContactCompany - get all contacts from company
func (e *Edb) GetContactCompany(id int64) ([]ContactCompany, error) {
	var contacts []ContactCompany
	_, err := e.db.Query(&contacts, `
		SELECT
			c.id,
			c.name,
			po.name AS post_name,
			pog.name AS post_go_name
		FROM
			contacts AS c
		LEFT JOIN
			posts AS po ON c.post_id = po.id
		LEFT JOIN
			posts AS pog ON c.post_go_id = pog.id
		WHERE
			c.company_id = $1
		ORDER BY
			name ASC
	`)
	if err != nil {
		log.Println("GetContactCompany e.db.Prepare ", err)
		return []ContactCompany{}, err
	}
	return contacts, err
}

// CreateContact - create new contact
func (e *Edb) CreateContact(contact Contact) (int64, error) {
	err := e.db.Insert(&contact)
	if err != nil {
		log.Println("CreateContact e.db.Insert ", err)
		return 0, err
	}
	_ = e.CreateContactEmails(contact)
	_ = e.CreateContactPhones(contact, false)
	_ = e.CreateContactPhones(contact, true)
	// CreateContactEducations(contact)
	return contact.ID, nil
}

// UpdateContact - save contact changes
func (e *Edb) UpdateContact(contact Contact) error {
	stmt, err := e.db.Prepare(`
		UPDATE
			contacts
		SET
			name=$2,
			company_id=$3,
			department_id=$4,
			post_id=$5,
			post_go_id=$6,
			rank_id=$7,
			birthday=$8,
			note=$9,
			updated_at = now()
		WHERE
			id = $1
	`)
	if err != nil {
		log.Println("UpdateContact e.db.Prepare ", err)
		return err
	}
	_, err = stmt.Exec(i2n(contact.ID), s2n(contact.Name), i2n(contact.CompanyID), i2n(contact.DepartmentID), i2n(contact.PostID), i2n(contact.PostGOID), i2n(contact.RankID), sd2n(contact.Birthday), s2n(contact.Note))
	if err != nil {
		log.Println("UpdateContact stmt.Exec ", err)
		return err
	}
	_ = e.CreateContactEmails(contact)
	_ = e.CreateContactPhones(contact, false)
	_ = e.CreateContactPhones(contact, true)
	// CreateContactEducations(contact)
	return nil
}

// DeleteContact - delete contact by id
func (e *Edb) DeleteContact(id int64) error {
	if id == 0 {
		return nil
	}
	err := e.DeleteAllContactPhones(id)
	if err != nil {
		log.Println("DeleteContact DeleteAllContactPhones ", err)
		return err
	}
	e.db.Exec(`
		DELETE FROM
			contacts
		WHERE
			id = $1
	`, id)
	if err != nil {
		log.Println("DeleteContact e.db.Exec ", id, err)
	}
	return err
}

func (e *Edb) contactCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			contacts (
				id bigserial primary key,
				name text,
				company_id bigint,
				department_id bigint,
				post_id bigint,
				post_go_id bigint,
				rank_id bigint,
				birthday date,
				note text,
				created_at TIMESTAMP without time zone,
				updated_at TIMESTAMP without time zone,
				UNIQUE(name, birthday)
			)
	`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("contactCreateTable ", err)
	}
	return err
}
