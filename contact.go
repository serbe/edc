package edc

import (
	"context"
	"time"
)

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

// ContactGet - get one contact by id
func ContactGet(id int64) (Contact, error) {
	var contact Contact
	if id == 0 {
		return contact, nil
	}
	contact.ID = id
	err := pool.QueryRow(context.Background(), `
		SELECT
			c.name,
			c.company_id,
			c.department_id,
			c.post_id,
			c.post_go_id,
			c.rank_id,
			c.birthday,
			c.note,
			c.created_at,
			c.updated_at,
			array_agg(DISTINCT e.email) AS emails,
			array_agg(DISTINCT ph.phone) AS phones,
			array_agg(DISTINCT f.phone) AS faxes,
			array_agg(DISTINCT ed.start_date) AS educations
		FROM
			contacts AS c
		LEFT JOIN
			emails AS e ON c.id = e.contact_id
		LEFT JOIN
			phones AS ph ON c.id = ph.contact_id AND ph.fax = false
		LEFT JOIN
			phones AS f ON c.id = f.contact_id AND f.fax = true
		LEFT JOIN
			educations AS ed ON c.id = ed.contact_id
		WHERE
			c.id = $1
		GROUP BY
			c.id
	`, id).Scan(&contact.Name, &contact.CompanyID, &contact.DepartmentID, &contact.PostID, &contact.PostGOID, &contact.RankID,
		&contact.Birthday, &contact.Note, &contact.CreatedAt, &contact.UpdatedAt, &contact.Emails, &contact.Phones, &contact.Faxes, &contact.Educations)
	if err != nil {
		errmsg("GetContact QueryRow", err)
		return contact, err
	}
	return contact, err
}

// ContactListGet - get all contacts for list
func ContactListGet() ([]ContactList, error) {
	var contacts []ContactList
	rows, err := pool.Query(context.Background(), `
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
		errmsg("GetContactList Query", err)
	}
	for rows.Next() {
		var contact ContactList
		err := rows.Scan(&contact.ID, &contact.Name, &contact.CompanyID, &contact.CompanyName,
			&contact.PostName, &contact.Phones, &contact.Faxes)
		if err != nil {
			errmsg("GetContactList Scan", err)
			return contacts, err
		}
		contacts = append(contacts, contact)
	}
	return contacts, rows.Err()
}

// ContactSelectGet - get all contacts for select
func ContactSelectGet() ([]SelectItem, error) {
	var contacts []SelectItem
	rows, err := pool.Query(context.Background(), `
		SELECT
			id,
			name
		FROM
			contacts
		ORDER BY
			name ASC
	`)
	if err != nil {
		errmsg("ContactSelectGet Query", err)
	}
	for rows.Next() {
		var contact SelectItem
		err := rows.Scan(&contact.ID, &contact.Name)
		if err != nil {
			errmsg("ContactSelectGet select", err)
			return contacts, err
		}
		contacts = append(contacts, contact)
	}
	return contacts, rows.Err()
}

// ContactCompanyGet - get all contacts from company
func ContactCompanyGet(id int64) ([]ContactShort, error) {
	var contacts []ContactShort
	if id == 0 {
		return contacts, nil
	}
	rows, err := pool.Query(context.Background(), `
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
	for rows.Next() {
		var contact ContactShort
		err := rows.Scan(&contact.ID, &contact.Name, &contact.PostName, &contact.PostGOName)
		if err != nil {
			errmsg("GetCompanyList Scan", err)
			return contacts, err
		}
		contacts = append(contacts, contact)
	}
	return contacts, rows.Err()
}

// ContactInsert - create new contact
func ContactInsert(contact Contact) (int64, error) {
	err := pool.QueryRow(context.Background(), `
		INSERT INTO contacts
		(
			name,
			company_id,
			department_id,
			post_id,
			post_go_id,
			rank_id,
			birthday,
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
			$7,
			$8,
			$9,
			$10
		)
		RETURNING
			id
	`, contact.Name, contact.CompanyID, contact.DepartmentID, contact.PostID, contact.PostGOID, contact.RankID, contact.Birthday, contact.Note,
		time.Now(), time.Now()).Scan(&contact.ID)
	if err != nil {
		errmsg("ContactInsert QueryRow", err)
		return 0, err
	}
	_ = EmailContactUpdate(contact.ID, contact.Emails)
	_ = PhoneContactUpdate(contact.ID, contact.Phones, false)
	_ = PhoneContactUpdate(contact.ID, contact.Faxes, true)
	return contact.ID, nil
}

// ContactUpdate - save contact changes
func ContactUpdate(contact Contact) error {
	_, err := pool.Exec(context.Background(), `
		UPDATE contacts SET
			name = $2,
			company_id = $3,
			department_id = $4,
			post_id = $5,
			post_go_id = $6,
			rank_id = $7,
			birthday = $8,
			note = $9,
			updated_at = $10
		WHERE
			id = $1
	`, contact.ID, contact.Name, contact.CompanyID, contact.DepartmentID, contact.PostID, contact.PostGOID, contact.RankID, contact.Birthday,
		contact.Note, time.Now())
	if err != nil {
		errmsg("ContactUpdate Exec", err)
		return err
	}
	_ = EmailContactUpdate(contact.ID, contact.Emails)
	_ = PhoneContactUpdate(contact.ID, contact.Phones, false)
	_ = PhoneContactUpdate(contact.ID, contact.Faxes, true)
	return nil
}

// ContactDelete - delete contact by id
func ContactDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_ = EmailContactDelete(id)
	_ = PhoneContactDelete(id, true)
	_ = PhoneContactDelete(id, false)
	_, err := pool.Exec(context.Background(), `
		DELETE FROM
			contacts
		WHERE
			id = $1
	`, id)
	if err != nil {
		errmsg("ContactDelete Exec", err)
	}
	return err
}

func contactCreateTable() error {
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
	_, err := pool.Exec(context.Background(), str)
	if err != nil {
		errmsg("contactCreateTable exec", err)
	}
	return err
}
