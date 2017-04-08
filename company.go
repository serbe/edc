package edc

// Company is struct for company
type Company struct {
	ID        int64          `sql:"id"             json:"id"`
	Name      string         `sql:"name"           json:"name"`
	Address   string         `sql:"address, null"  json:"address"`
	Scope     SelectItem     `sql:"-"              json:"scope"`
	ScopeID   int64          `sql:"scope_id, null" json:"scope_id"`
	Note      string         `sql:"note, null"     json:"note"`
	Emails    []Email        `sql:"-"              json:"emails"`
	Phones    []Phone        `sql:"-"              json:"phones"`
	Faxes     []Phone        `sql:"-"              json:"faxes"`
	Practices []PracticeList `sql:"-"              json:"practices"`
	Contacts  []ContactTiny  `sql:"-"              json:"contacts"`
	CreatedAt string         `sql:"created_at"     json:"-"`
	UpdatedAt string         `sql:"updated_at"     json:"-"`
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
	err := e.db.Model(&company).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetCompany select", err)
		return company, err
	}
	if company.ScopeID > 0 {
		company.Scope, err = e.GetScopeSelect(company.ScopeID)
		if err != nil {
			errmsg("GetCompany e.GetScope ", err)
			return company, err
		}
	}
	company.Emails, err = e.GetCompanyEmails(company.ID)
	if err != nil {
		errmsg("GetCompany e.GetCompanyEmails ", err)
		return company, err
	}
	company.Phones, err = e.GetCompanyPhones(company.ID, false)
	if err != nil {
		errmsg("GetCompany e.GetCompanyPhones false ", err)
		return company, err
	}
	company.Faxes, err = e.GetCompanyPhones(company.ID, true)
	if err != nil {
		errmsg("GetCompany e.GetCompanyPhones true ", err)
		return company, err
	}
	company.Practices, err = e.GetPracticeCompany(id)
	if err != nil {
		errmsg("GetCompany e.GetPracticeCompany", err)
	}
	company.Contacts, err = e.GetContactCompany(company.ID)
	if err != nil {
		errmsg("GetCompany e.GetContactCompany ", err)
		return company, err
	}
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
			array_agg(DISTINCT e.email) AS emails,
			array_agg(DISTINCT p.phone) AS phones,
			array_agg(DISTINCT f.phone) AS faxes,
			array_agg(DISTINCT pr.date_of_practice) AS practices
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
		errmsg("GetCompanyList query", err)
	}
	return companies, err
}

// GetCompanySelect - get company for contact
func (e *Edb) GetCompanySelect(id int64) (SelectItem, error) {
	var company SelectItem
	err := e.db.Model(&Company{}).
		Column("id", "name").
		Where("id = ?", id).
		Select(&company)
	if err != nil {
		errmsg("GetCompanySelect select", err)
	}
	return company, err
}

// GetCompanySelectAll - get all companyes for select
func (e *Edb) GetCompanySelectAll() ([]SelectItem, error) {
	var companies []SelectItem
	err := e.db.Model(&Company{}).
		Column("id", "name").
		Order("name ASC").
		Select(&companies)
	if err != nil {
		errmsg("GetCompanySelectAll select", err)
	}
	return companies, err
}

// CreateCompany - create new company
func (e *Edb) CreateCompany(company Company) (int64, error) {
	err := e.db.Insert(&company)
	if err != nil {
		errmsg("CreateCompany insert", err)
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
		errmsg("UpdateCompany update", err)
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
	err := e.DeleteAllCompanyPhones(id)
	if err != nil {
		errmsg("DeleteCompany DeleteAllCompanyPhones", err)
	}
	_, err = e.db.Model(&Company{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteCompany delete", err)
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
				updated_at timestamp without time zone default now(),
				UNIQUE(name, scope_id)
			)
	`
	_, err := e.db.Exec(str)
	if err != nil {
		errmsg("companyCreateTable exec", err)
	}
	return err
}
