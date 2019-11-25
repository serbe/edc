package edc

// Contact is struct for contact
type Contact struct {
	ID           int64    `sql:"id"            json:"id"            form:"id"            query:"id"`
	Name         string   `sql:"name"          json:"name"          form:"name"          query:"name"`
	CompanyID    int64    `sql:"company_id"    json:"company_id"    form:"company_id"    query:"company_id"`
	DepartmentID int64    `sql:"department_id" json:"department_id" form:"department_id" query:"department_id"`
	PostID       int64    `sql:"post_id"       json:"post_id"       form:"post_id"       query:"post_id"`
	PostGOID     int64    `sql:"post_go_id"    json:"post_go_id"    form:"post_go_id"    query:"post_go_id"`
	RankID       int64    `sql:"rank_id"       json:"rank_id"       form:"rank_id"       query:"rank_id"`
	Birthday     string   `sql:"birthday"      json:"birthday"      form:"birthday"      query:"birthday"`
	Note         string   `sql:"note"          json:"note"          form:"note"          query:"note"`
	CreatedAt    string   `sql:"created_at"    json:"-"`
	UpdatedAt    string   `sql:"updated_at"    json:"-"`
	Emails       []string `sql:"-"             json:"emails"        form:"emails"        query:"emails"`
	Phones       []int64  `sql:"-"             json:"phones"        form:"phones"        query:"phones"`
	Faxes        []int64  `sql:"-"             json:"faxes"         form:"faxes"         query:"faxes"`
	Educations   []string `sql:"-"             json:"educations"    form:"educations"    query:"educations"`
}

// ContactList is struct for contact list
type ContactList struct {
	ID          int64   `json:"id"           form:"id"           query:"id"`
	Name        string  `json:"name"         form:"name"         query:"name"`
	CompanyID   int64   `json:"company_id"   form:"company_id"   query:"company_id"`
	CompanyName string  `json:"company_name" form:"company_name" query:"company_name"`
	PostName    string  `json:"post_name"    form:"post_name"    query:"post_name"`
	Phones      []int64 `json:"phones"       form:"phones"       query:"phones"       pg:",array"`
	Faxes       []int64 `json:"faxes"        form:"faxes"        query:"faxes"        pg:",array"`
}

// ContactShort is struct of contact for another parents
type ContactShort struct {
	ID             int64  `json:"id"              form:"id"              query:"id"`
	Name           string `json:"name"            form:"name"            query:"name"`
	DepartmentName string `json:"department_name" form:"department_name" query:"department_name"`
	PostName       string `json:"post_name"       form:"post_name"       query:"post_name"`
	PostGOName     string `json:"post_go_name"    form:"post_go_name"    query:"post_go_name"`
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
		return contact, err
	}
	// contact.Educations, err = e.ContactEducations(contact.ID)
	// if err != nil {
	// 	errmsg("GetContact ContactEducations", err)
	// 	return Contact{}, err
	// }
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

// GetContactSelect - get contact for select by id
func (e *Edb) GetContactSelect(id int64) (SelectItem, error) {
	var contact SelectItem
	err := e.db.Model(&Contact{}).
		Column("id", "name").
		Where("id = ?", id).
		Select(&contact)
	if err != nil {
		errmsg("GetContactSelect select", err)
	}
	return contact, err
}

// GetContactSelectAll - get all contacts for select
func (e *Edb) GetContactSelectAll() ([]SelectItem, error) {
	var contacts []SelectItem
	err := e.db.Model(&Contact{}).
		Column("id", "name").
		Order("name ASC").
		Select(&contacts)
	if err != nil {
		errmsg("GetContactSelectAll select", err)
	}
	return contacts, err
}

// GetContactCompany - get all contacts from company
func (e *Edb) GetContactCompany(id int64) ([]ContactShort, error) {
	var contacts []ContactShort
	if id == 0 {
		return contacts, nil
	}
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
	_ = e.UpdateContactEmails(contact)
	_ = e.UpdateContactPhones(contact, false)
	_ = e.UpdateContactPhones(contact, true)
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
	_ = e.UpdateContactEmails(contact)
	_ = e.UpdateContactPhones(contact, false)
	_ = e.UpdateContactPhones(contact, true)
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
