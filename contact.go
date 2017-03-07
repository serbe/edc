package edc

// Contact is struct for contact
type Contact struct {
	ID           int64       `sql:"id"                  json:"id"`
	Name         string      `sql:"name"                json:"name"`
	Company      Company     `sql:"-"                   json:"company"`
	CompanyID    int64       `sql:"company_id, null"    json:"company_id"`
	Department   Department  `sql:"-"                   json:"department"`
	DepartmentID int64       `sql:"department_id, null" json:"department_id"`
	Post         Post        `sql:"-"                   json:"post"`
	PostID       int64       `sql:"post_id, null"       json:"post_id"`
	PostGO       Post        `sql:"-"                   json:"post_go"`
	PostGOID     int64       `sql:"post_go_id, null"    json:"post_go_id"`
	Rank         Rank        `sql:"-"                   json:"rank"`
	RankID       int64       `sql:"rank_id, null"       json:"rank_id"`
	Birthday     string      `sql:"birthday, null"      json:"birthday"`
	Note         string      `sql:"note, null"          json:"note"`
	Emails       []Email     `sql:"-"                   json:"emails"`
	Phones       []Phone     `sql:"-"                   json:"phones"`
	Faxes        []Phone     `sql:"-"                   json:"faxes"`
	Educations   []Education `sql:"-"                   json:"educations"`
	CreatedAt    string      `sql:"created_at"          json:"created_at"`
	UpdatedAt    string      `sql:"updated_at"          json:"updated_at"`
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
	var contact Contact
	if id == 0 {
		return contact, nil
	}
	err := e.db.Model(&contact).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetContact select", err)
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
			array_agg(DISTINCT ph.phone) AS phones,
			array_agg(DISTINCT f.phone) AS faxes
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
		errmsg("GetContactList query", err)
	}
	return contacts, err
}

// GetContactSelect - get all contacts for select
func (e *Edb) GetContactSelect() ([]SelectItem, error) {
	var contacts []SelectItem
	err := e.db.Model(&Contact{}).
		Column("id", "name").
		Order("name ASC").
		Select(&contacts)
	if err != nil {
		errmsg("GetContactSelect select", err)
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
			c.company_id = ?
		ORDER BY
			name ASC
	`, id)
	if err != nil {
		errmsg("GetContactCompany query", err)
	}
	return contacts, err
}

// CreateContact - create new contact
func (e *Edb) CreateContact(contact Contact) (int64, error) {
	err := e.db.Insert(&contact)
	if err != nil {
		errmsg("CreateContact insert", err)
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
	err := e.db.Update(&contact)
	if err != nil {
		errmsg("UpdateContact update", err)
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
		errmsg("DeleteContact DeleteAllContactPhones", err)
		return err
	}
	_, err = e.db.Model(&Contact{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteContact delete", err)
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
				updated_at TIMESTAMP without time zone default now(),
				UNIQUE(name, birthday)
			)
	`
	_, err := e.db.Exec(str)
	if err != nil {
		errmsg("contactCreateTable exec", err)
	}
	return err
}
