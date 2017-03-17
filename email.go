package edc

// Email - struct for email
type Email struct {
	ID        int64  `sql:"id"                   json:"id"`
	CompanyID int64  `sql:"company_id, pk, null" json:"company_id"`
	ContactID int64  `sql:"contact_id, pk, null" json:"contact_id"`
	Email     string `sql:"email, null"          json:"email"`
	CreatedAt string `sql:"created_at"           json:"created_at"`
	UpdatedAt string `sql:"updated_at"           json:"updated_at"`
}

// GetEmail - get one email by id
func (e *Edb) GetEmail(id int64) (Email, error) {
	var email Email
	if id == 0 {
		return email, nil
	}
	err := e.db.Model(&email).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetEmail select", err)
	}
	return email, nil
}

// GetEmails - get all emails for list
func (e *Edb) GetEmails() ([]Email, error) {
	var emails []Email
	err := e.db.Model(&emails).
		Column("id", "email").
		Order("email ASC").
		Select()
	if err != nil {
		errmsg("GetEmailList select", err)
	}
	return emails, err
}

// GetCompanyEmails - get all emails by company id
func (e *Edb) GetCompanyEmails(id int64) ([]Email, error) {
	var emails []Email
	if id == 0 {
		return emails, nil
	}
	err := e.db.Model(&emails).
		Column("id", "email").
		Order("email ASC").
		Where("company_id = ?", id).Select()
	if err != nil {
		errmsg("GetCompanyEmails select", err)
	}
	return emails, err
}

// GetContactEmails - get all emails by contact id
func (e *Edb) GetContactEmails(id int64) ([]Email, error) {
	var emails []Email
	if id == 0 {
		return emails, nil
	}
	err := e.db.Model(&emails).
		Column("id", "email").
		Order("email ASC").
		Where("contact_id = ?", id).Select()
	if err != nil {
		errmsg("GetContactEmails select", err)
	}
	return emails, err
}

// CreateEmail - create new email
func (e *Edb) CreateEmail(email Email) (int64, error) {
	err := e.db.Insert(&email)
	if err != nil {
		errmsg("CreateEmail insert", err)
	}
	return email.ID, nil
}

// CreateCompanyEmails - create new company email
func (e *Edb) CreateCompanyEmails(company Company) error {
	err := e.DeleteCompanyEmails(company.ID)
	if err != nil {
		errmsg("CreateCompanyEmails DeleteCompanyEmails", err)
		return err
	}
	for _, email := range company.Emails {
		email.CompanyID = company.ID
		_, err = e.CreateEmail(email)
		if err != nil {
			errmsg("CreateCompanyEmails CreateEmail", err)
			return err
		}
	}
	return nil
}

// CreateContactEmails - create new contact email
func (e *Edb) CreateContactEmails(contact Contact) error {
	err := e.DeleteContactEmails(contact.ID)
	if err != nil {
		errmsg("CreateContactEmails DeleteContactEmails", err)
		return err
	}
	for _, email := range contact.Emails {
		email.ContactID = contact.ID
		_, err = e.CreateEmail(email)
		if err != nil {
			errmsg("CreateContactEmails CreateEmail", err)
			return err
		}
	}
	return nil
}

// UpdateEmail - save email changes
func (e *Edb) UpdateEmail(email Email) error {
	err := e.db.Update(&email)
	if err != nil {
		errmsg("UpdateEmail update", err)
	}
	return err
}

// DeleteEmail - delete email by id
func (e *Edb) DeleteEmail(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Model(&Email{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteEmail delete", err)
	}
	return err
}

// DeleteCompanyEmails - delete all emails by company id
func (e *Edb) DeleteCompanyEmails(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Model(&Email{}).
		Where("company_id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteCompanyEmails delete", err)
	}
	return err
}

// DeleteContactEmails - delete all emails by contact id
func (e *Edb) DeleteContactEmails(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Model(&Email{}).
		Where("contact_id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteContactEmails delete", err)
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
				updated_at timestamp without time zone default now()
			)
	`
	_, err := e.db.Exec(str)
	if err != nil {
		errmsg("emailCreateTable exec", err)
	}
	return err
}
