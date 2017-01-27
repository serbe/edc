package edc

import "log"

// Company is struct for company
type Company struct {
	ID        int64            `sql:"id" json:"id"`
	Name      string           `sql:"name" json:"name"`
	Address   string           `sql:"address, null" json:"address"`
	Scope     Scope            `sql:"-"`
	ScopeID   int64            `sql:"scope_id, null" json:"scope_id"`
	Note      string           `sql:"note, null" json:"note"`
	Emails    []Email          `sql:"-"`
	Phones    []Phone          `sql:"-"`
	Faxes     []Phone          `sql:"-"`
	Practices []Practice       `sql:"-"`
	Contacts  []ContactCompany `sql:"-"`
	CreatedAt string           `sql:"created_at" json:"created_at"`
	UpdatedAt string           `sql:"updated_at" json:"updated_at"`
}

// CompanyList is struct for list company
type CompanyList struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	Address   string   `json:"address"`
	ScopeName string   `json:"scope_name"`
	Emails    []string `json:"emails"      pg:",array"`
	Phones    []string `json:"phones"      pg:",array"`
	Faxes     []string `json:"faxes"       pg:",array"`
	Practices []string `json:"practices"   pg:",array"`
}

// GetCompany - get one company by id
func (e *Edb) GetCompany(id int64) (Company, error) {
	var company Company
	if id == 0 {
		return company, nil
	}
	err := e.db.Model(&company).Where("id = ?", id).Select()
	company.Practices, err = e.GetPracticeCompany(id)
	return company, err
}

// GetCompanyList - get all companyes for list
func (e *Edb) GetCompanyList() ([]CompanyList, error) {
	var companies []CompanyList
	_, err := e.db.Query(&companies, `
		SELECT
			c.id,
			c.name,
			c.address,
			s.name AS scope_name,
			array_to_string(array_agg(DISTINCT e.email),',') AS email,
			array_to_string(array_agg(DISTINCT p.phone),',') AS phone,
			array_to_string(array_agg(DISTINCT f.phone),',') AS fax,
			array_to_string(array_agg(DISTINCT pr.date_of_practice),',') AS practice
        FROM
			companies AS c
		LEFT JOIN
			scopes AS s ON c.scope_id = s.id
		LEFT JOIN
			emails AS e ON c.id = e.company_id
		LEFT JOIN
			phones AS p ON c.id = p.company_id AND p.fax = false
		LEFT JOIN
			phones AS f ON c.id = f.company_id AND f.fax = true
		LEFT JOIN
			practices AS pr ON c.id = pr.company_id
		GROUP BY
			c.id,
			s.name
		ORDER BY
			c.name ASC
	`)
	if err != nil {
		log.Println("GetCompanyList e.db.Query ", err)
		return []CompanyList{}, err
	}
	return companies, err
}

// GetCompanySelect - get all companyes for select
func (e *Edb) GetCompanySelect() ([]SelectItem, error) {
	var companies []SelectItem
	rows, err := e.db.Query(&companies, `
		SELECT
			c.id,
			c.name
        FROM
			companies AS c
		ORDER BY
			c.name ASC
	`)
	if err != nil {
		log.Println("GetCompanyList e.db.Query ", err)
		return []SelectItem{}, err
	}
	return companies, err
}

// CreateCompany - create new company
func (e *Edb) CreateCompany(company Company) (int64, error) {
	err := e.db.Insert(&company)
	if err != nil {
		log.Println("CreateCompany e.db.Insert ", err)
		return 0, err
	}
	_ = e.CreateCompanyEmails(company)
	_ = e.CreateCompanyPhones(company, false)
	_ = e.CreateCompanyPhones(company, true)
	return company.ID, nil
}

// UpdateCompany - save company changes
func (e *Edb) UpdateCompany(company Company) error {
	err := e.db.Update(&company)
	if err != nil {
		log.Println("UpdateCompany e.db.Update ", err)
		return err
	}
	_ = e.CreateCompanyEmails(company)
	_ = e.CreateCompanyPhones(company, false)
	_ = e.CreateCompanyPhones(company, true)
	return nil
}

// DeleteCompany - delete company by id
func (e *Edb) DeleteCompany(id int64) error {
	if id == 0 {
		return nil
	}
	c, err := e.GetCompany(id)
	if err != nil {
		log.Println("DeleteCompany e.GetCompany ", err)
		return err
	}
	e.DeleteAllCompanyPhones(id)
	err = e.db.Delete(&c)
	if err != nil {
		log.Println("DeleteCompany e.db.Delete ", id, err)
	}
	return err
}

func (e *Edb) companyCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			companies (
				id BIGSERIAL PRIMARY KEY,
				name TEXT,
				address TEXT,
				scope_id BIGINT,
				note TEXT,
				created_at timestamp without time zone,
				updated_at timestamp without time zone,
				UNIQUE(name, scope_id)
			)
	`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("companyCreateTable e.db.Exec ", err)
	}
	return err
}
