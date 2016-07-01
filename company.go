package edc

import (
	"log"
)

// Company is struct for company
type Company struct {
	ID        int64      `sql:"id" json:"id"`
	Name      string     `sql:"name" json:"name"`
	Address   string     `sql:"address, null" json:"address"`
	Scope     Scope      `sql:"-"`
	ScopeID   int64      `sql:"scope_id, null" json:"scope-id"`
	Notes     string     `sql:"notes, null" json:"notes"`
	Emails    []Email    `sql:"-"`
	Phones    []Phone    `sql:"-"`
	Faxes     []Phone    `sql:"-"`
	Practices []Practice `sql:"-"`
}

// GetCompany - get one company by id
func (e *EDc) GetCompany(id int64) (company Company, err error) {
	if id == 0 {
		return
	}
	_, err = e.db.QueryOne(&company, "SELECT * FROM companies WHERE id = ? LIMIT 1", id)
	if err != nil {
		log.Println("GetCompany e.db.QueryRow Scan ", err)
		return
	}
	company.Scope, _ = e.GetScope(company.ScopeID)
	company.Emails, _ = e.GetCompanyEmails(company.ID)
	company.Phones, _ = e.GetCompanyPhones(company.ID)
	company.Faxes, _ = e.GetCompanyFaxes(company.ID)
	company.Practices, _ = e.GetCompanyPractices(company.ID)
	return
}

// GetCompanyAll - get all companyes
func (e *EDc) GetCompanyAll() (companyes []Company, err error) {
	_, err = e.db.Query(&companyes, "SELECT * FROM companies")
	if err != nil {
		log.Println("GetCompanyAll e.db.Query ", err)
		return
	}
	for i := range companyes {
		companyes[i].Scope, _ = e.GetScope(companyes[i].ScopeID)
		companyes[i].Emails, _ = e.GetCompanyEmails(companyes[i].ID)
		companyes[i].Phones, _ = e.GetCompanyPhones(companyes[i].ID)
		companyes[i].Faxes, _ = e.GetCompanyFaxes(companyes[i].ID)
		companyes[i].Practices, _ = e.GetCompanyPractices(companyes[i].ID)
		for j := range companyes[i].Practices {
			companyes[i].Practices[j].DateStr = companyes[i].Practices[j].DateOfPractice.Format("02.01.2006")
		}
	}
	return
}

// CreateCompany - create new company
func (e *EDc) CreateCompany(company Company) (err error) {
	err = e.db.Create(&company)
	if err != nil {
		log.Println("CreateCompany e.db.Create ", err)
		return
	}
	e.CreateCompanyEmails(company)
	e.CreateCompanyPhones(company)
	e.CreateCompanyFaxes(company)
	return
}

// UpdateCompany - save company changes
func (e *EDc) UpdateCompany(company Company) (err error) {
	err = e.db.Update(&company)
	if err != nil {
		log.Println("UpdateCompany e.db.Update ", err)
		return
	}
	e.CreateCompanyEmails(company)
	e.CreateCompanyPhones(company)
	e.CreateCompanyFaxes(company)
	// CreateCompanyPractices(c)
	return
}

// DeleteCompany - delete company by id
func (e *EDc) DeleteCompany(id int64) (err error) {
	if id == 0 {
		return
	}
	e.DeleteAllCompanyPhones(id)
	_, err = e.db.Exec("DELETE FROM companies WHERE id=?", id)
	if err != nil {
		log.Println("DeleteCompany e.db.Exec ", err)
	}
	return
}

func (e *EDc) companyCreateTable() (err error) {
	str := `CREATE TABLE IF NOT EXISTS companies (id bigserial primary key, name text, address text, scope_id bigint, notes text)`
	_, err = e.db.Exec(str)
	if err != nil {
		log.Println("companyCreateTable e.db.Exec ", err)
	}
	return
}
