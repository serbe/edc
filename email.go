package edc

import "log"

// Email - struct for email
type Email struct {
	TableName struct{} `sql:"emails"`
	ID        int64    `sql:"id" json:"id"`
	CompanyID int64    `sql:"company_id, pk, null" json:"company-id"`
	PeopleID  int64    `sql:"people_id, pk, null" json:"people-id"`
	Email     string   `sql:"email, null" json:"email"`
	Notes     string   `sql:"notes, null" json:"notes"`
}

// GetEmail - get one email by id
func (e *EDc) GetEmail(id int64) (email Email, err error) {
	if id == 0 {
		return
	}
	err = e.db.Model(&email).Where("id = ?", id).Select()
	if err != nil {
		log.Println("GetEmail ", err)
	}
	return
}

// GetEmailAll - get all emails
func (e *EDc) GetEmailAll() (emails []Email, err error) {
	err = e.db.Model(&emails).Order("name ASC").Select()
	if err != nil {
		log.Println("GetEmailAll ", err)
		return
	}
	return
}

// GetCompanyEmails - get all emails by company id
func (e *EDc) GetCompanyEmails(id int64) (emails []Email, err error) {
	if id == 0 {
		return
	}
	err = e.db.Model(&emails).Where("company_id = ?", id).Order("name ASC").Select()
	if err != nil {
		log.Println("GetCompanyEmails ", err)
		return
	}
	return
}

// GetPeopleEmails - get all emails by people id
func (e *EDc) GetPeopleEmails(id int64) (emails []Email, err error) {
	if id == 0 {
		return
	}
	err = e.db.Model(&emails).Where("people_id = ?", id).Order("name ASC").Select()
	if err != nil {
		log.Println("GetPeopleEmails ", err)
		return
	}
	return
}

// CreateEmail - create new email
func (e *EDc) CreateEmail(email Email) (err error) {
	err = e.db.Create(&email)
	if err != nil {
		log.Println("CreateEmail ", err)
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
		log.Println("UpdateEmail ", err)
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
		log.Println("DeleteEmail ", err)
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
		log.Println("DeleteCompanyEmails ", err)
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
		log.Println("DeletePeopleEmails ", err)
	}
	return
}

func (e *EDc) emailCreateTable() (err error) {
	str := `CREATE TABLE IF NOT EXISTS emails (id bigserial primary key, company_id bigint, people_id bigint, email text, notes text)`
	_, err = e.db.Exec(str)
	if err != nil {
		log.Println("emailCreateTable ", err)
	}
	return
}
