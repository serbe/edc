package edc

import "log"

// Email - struct for email
type Email struct {
	ID        int64  `sql:"id" json:"id"`
	CompanyID int64  `sql:"company_id, pk, null" json:"company-id"`
	PeopleID  int64  `sql:"people_id, pk, null" json:"people-id"`
	Email     string `sql:"email, null" json:"email"`
	Notes     string `sql:"notes, null" json:"notes"`
}

// GetEmail - get one email by id
func (e *EDc) GetEmail(id int64) (email Email, err error) {
	if id == 0 {
		return
	}
	_, err = e.db.QueryOne(&email, "SELECT * FROM email WHERE id = ? LIMIT 1", id)
	if err != nil {
		log.Println("GetEmail e.db.QueryRow Scan ", err)
	}
	return
}

// GetEmailAll - get all emails
func (e *EDc) GetEmailAll() (emails []Email, err error) {
	_, err = e.db.Query(&emails, "SELECT * FROM emails")
	if err != nil {
		log.Println("GetEmailAll e.db.Query ", err)
		return
	}
	return
}

// GetCompanyEmails - get all emails by company id
func (e *EDc) GetCompanyEmails(id int64) (emails []Email, err error) {
	if id == 0 {
		return
	}
	_, err = e.db.Query(&emails, "SELECT * FROM emails WHERE company_id = ?", id)
	if err != nil {
		log.Println("GetCompanyEmails e.db.Query ", err)
		return
	}
	return
}

// GetPeopleEmails - get all emails by people id
func (e *EDc) GetPeopleEmails(id int64) (emails []Email, err error) {
	if id == 0 {
		return
	}
	_, err = e.db.Query(&emails, "SELECT * FROM emails WHERE people_id = ?", id)
	if err != nil {
		log.Println("GetPeopleEmails e.db.Query ", err)
		return
	}
	return
}

// CreateEmail - create new email
func (e *EDc) CreateEmail(email Email) (err error) {
	err = e.db.Create(&email)
	if err != nil {
		log.Println("CreateEmail e.db.Exec ", err)
	}
	return
}

// CreateCompanyEmails - create new company email
func (e *EDc) CreateCompanyEmails(company Company) (err error) {
	err = e.DeleteCompanyEmails(company.ID)
	if err != nil {
		log.Println("CreateCompanyEmails DeleteCompanyEmails ", err)
		return
	}
	for _, email := range company.Emails {
		email.CompanyID = company.ID
		err = e.CreateEmail(email)
		if err != nil {
			log.Println("CreateCompanyEmails CreateEmail ", err)
			return
		}
	}
	return
}

// CreatePeopleEmails - create new people email
func (e *EDc) CreatePeopleEmails(people People) (err error) {
	err = e.DeletePeopleEmails(people.ID)
	if err != nil {
		log.Println("CreatePeopleEmails DeletePeopleEmails ", err)
		return
	}
	for _, email := range people.Emails {
		email.PeopleID = people.ID
		err = e.CreateEmail(email)
		if err != nil {
			log.Println("CreatePeopleEmails CreateEmail ", err)
			return
		}
	}
	return
}

// UpdateEmail - save email changes
func (e *EDc) UpdateEmail(email Email) (err error) {
	e.db.Update(&email)
	if err != nil {
		log.Println("UpdateEmail e.db.Update ", err)
	}
	return
}

// DeleteEmail - delete email by id
func (e *EDc) DeleteEmail(id int64) (err error) {
	if id == 0 {
		return
	}
	_, err = e.db.Exec("DELETE FROM emails WHERE id = ?", id)
	if err != nil {
		log.Println("DeleteEmail e.db.Exec ", err)
	}
	return
}

// DeleteCompanyEmails - delete all emails by company id
func (e *EDc) DeleteCompanyEmails(id int64) (err error) {
	if id == 0 {
		return
	}
	_, err = e.db.Exec("DELETE FROM emails WHERE company_id = ?", id)
	if err != nil {
		log.Println("DeleteCompanyEmails e.db.Exec ", err)
	}
	return
}

// DeletePeopleEmails - delete all emails by people id
func (e *EDc) DeletePeopleEmails(id int64) (err error) {
	if id == 0 {
		return
	}
	_, err = e.db.Exec("DELETE FROM emails WHERE people_id = ?", id)
	if err != nil {
		log.Println("DeletePeopleEmails e.db.Exec ", err)
	}
	return
}

func (e *EDc) emailCreateTable() (err error) {
	str := `CREATE TABLE IF NOT EXISTS emails (id bigserial primary key, company_id bigint, people_id bigint, email text, notes text)`
	_, err = e.db.Exec(str)
	if err != nil {
		log.Println("emailCreateTable e.db.Exec ", err)
	}
	return
}
