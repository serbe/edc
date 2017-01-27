package edc

import "log"

// Email - struct for email
type Email struct {
	ID        int64  `sql:"id" json:"id"`
	CompanyID int64  `sql:"company_id, pk, null" json:"company_id"`
	ContactID int64  `sql:"contact_id, pk, null" json:"contact_id"`
	Email     string `sql:"email, null" json:"email"`
	CreatedAt string `sql:"created_at" json:"created_at"`
	UpdatedAt string `sql:"updated_at" json:"updated_at"`
}

// GetEmail - get one email by id
func (e *Edb) GetEmail(id int64) (Email, error) {
	var email Email
	if id == 0 {
		return email, nil
	}
	err := e.db.Model(&email).Where(`id = ?`, id).Select()
	if err != nil {
		log.Println("GetEmail e.db.Select", err)
		return Email{}, err
	}
	return email, nil
}

// GetEmails - get all emails for list
func (e *Edb) GetEmails() ([]Email, error) {
	var emails []Email
	_, err := e.db.Query(&emails, `
		SELECT
			id,
			email
		FROM
			emails
		ORDER BY
			name ASC
	`)
	if err != nil {
		log.Println("GetEmailList e.db.Query ", err)
		return []Email{}, err
	}
	return emails, err
}

// GetCompanyEmails - get all emails by company id
func (e *Edb) GetCompanyEmails(id int64) ([]Email, error) {
	var emails []Email
	if id == 0 {
		return emails, nil
	}
	_, err := e.db.Query(&emails, `
		SELECT
			id,
			email
		FROM
			emails
		WHERE
			company_id = $1
		ORDER BY
			name ASC
	`, id)
	if err != nil {
		log.Println("GetCompanyEmails e.db.Query ", err)
		return []Email{}, err
	}
	return emails, err
}

// GetContactEmails - get all emails by contact id
func (e *Edb) GetContactEmails(id int64) ([]Email, error) {
	var emails []Email
	if id == 0 {
		return emails, nil
	}
	rows, err := e.db.Query(&emails, `
		SELECT
			id,
			email
		FROM
			emails
		WHERE
			contact_id = $1
		ORDER BY
			name ASC
	`, id)
	if err != nil {
		log.Println("GetContactEmails e.db.Query ", err)
		return []Email{}, err
	}
	return emails, err
}

// CreateEmail - create new email
func (e *Edb) CreateEmail(email Email) (int64, error) {
	err := e.db.Insert(&email)
	if err != nil {
		log.Println("CreateEmail e.db.Insert ", err)
		return 0, err
	}
	return email.ID, nil
}

// CreateCompanyEmails - create new company email
func (e *Edb) CreateCompanyEmails(company Company) error {
	err := e.DeleteCompanyEmails(company.ID)
	if err != nil {
		log.Println("CreateCompanyEmails DeleteCompanyEmails ", err)
		return err
	}
	for _, email := range company.Emails {
		email.CompanyID = company.ID
		_, err = e.CreateEmail(email)
		if err != nil {
			log.Println("CreateCompanyEmails CreateEmail ", err)
			return err
		}
	}
	return nil
}

// CreateContactEmails - create new contact email
func (e *Edb) CreateContactEmails(contact Contact) error {
	err := e.DeleteContactEmails(contact.ID)
	if err != nil {
		log.Println("CreateContactEmails DeleteContactEmails ", err)
		return err
	}
	for _, email := range contact.Emails {
		email.ContactID = contact.ID
		_, err = e.CreateEmail(email)
		if err != nil {
			log.Println("CreateContactEmails CreateEmail ", err)
			return err
		}
	}
	return nil
}

// UpdateEmail - save email changes
func (e *Edb) UpdateEmail(email Email) error {
	err := e.db.Update(&email)
	if err != nil {
		log.Println("UpdateEmail e.db.Update ", err)
		return err
	}
	return err
}

// DeleteEmail - delete email by id
func (e *Edb) DeleteEmail(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec(`
		DELETE FROM
			emails
		WHERE
			id = $1
	`, id)
	if err != nil {
		log.Println("DeleteEmail ", err)
	}
	return err
}

// DeleteCompanyEmails - delete all emails by company id
func (e *Edb) DeleteCompanyEmails(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec(`
		DELETE FROM
			emails
		WHERE
			company_id = $1
	`, id)
	if err != nil {
		log.Println("DeleteCompanyEmails ", id, err)
	}
	return err
}

// DeleteContactEmails - delete all emails by contact id
func (e *Edb) DeleteContactEmails(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec(`
		DELETE FROM
			emails
		WHERE
			contact_id = $1
	`, id)
	if err != nil {
		log.Println("DeleteContactEmails ", err)
	}
	return err
}

func (e *Edb) emailCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			emails (
				id bigserial primary key,
				company_id bigint,
				contact_id bigint,
				email text,
				created_at timestamp without time zone,
				updated_at timestamp without time zone
			)
	`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("emailCreateTable ", err)
	}
	return err
}
