package edc

import "context"

// Company is struct for company
type Company struct {
	ID        int64          `sql:"id"         json:"id"        form:"id"        query:"id"`
	Name      string         `sql:"name"       json:"name"      form:"name"      query:"name"`
	Address   string         `sql:"address"    json:"address"   form:"address"   query:"address"`
	ScopeID   int64          `sql:"scope_id"   json:"scope_id"  form:"scope_id"  query:"scope_id"`
	Note      string         `sql:"note"       json:"note"      form:"note"      query:"note"`
	CreatedAt string         `sql:"created_at" json:"-"`
	UpdatedAt string         `sql:"updated_at" json:"-"`
	Emails    []string       `sql:"-"          json:"emails"    form:"emails"    query:"emails"`
	Phones    []int64        `sql:"-"          json:"phones"    form:"phones"    query:"phones"`
	Faxes     []int64        `sql:"-"          json:"faxes"     form:"faxes"     query:"faxes"`
	Practices []PracticeList `sql:"-"          json:"practices" form:"practices" query:"practices"`
	Contacts  []ContactShort `sql:"-"          json:"contacts"  form:"contacts"  query:"contacts"`
}

// CompanyList is struct for list company
type CompanyList struct {
	ID        int64    `json:"id"         form:"id"         query:"id"`
	Name      string   `json:"name"       form:"name"       query:"name"`
	Address   string   `json:"address"    form:"address"    query:"address"`
	ScopeName string   `json:"scope_name" form:"scope_name" query:"scope_name"`
	Emails    []string `json:"emails"     form:"emails"     query:"emails"      pg:",array"`
	Phones    []int64  `json:"phones"     form:"phones"     query:"phones"      pg:",array"`
	Faxes     []int64  `json:"faxes"      form:"faxes"      query:"faxes"       pg:",array"`
	Practices []string `json:"practices"  form:"practices"  query:"practices"   pg:",array"`
}

// CompanyGet - get one company by id
func CompanyGet(id int64) (Company, error) {
	var company Company
	if id == 0 {
		return company, nil
	}
	var (
		emails    []string
		phones    []int64
		faxes     []int64
		practices []PracticeList
		contacts  []ContactShort
	)
	company.ID = id
	err := pool.QueryRow(context.Background(), `
		SELECT
			c.name,
			c.address,
			c.scope_id,
			c.note,
			c.created_at,
			c.updated_at,
			array_agg(DISTINCT e.email) AS emails,
			array_agg(DISTINCT ph.phone) AS phones,
			array_agg(DISTINCT f.phone) AS faxes
		FROM
			companies AS c
		LEFT JOIN
			emails AS e ON c.id = e.company_id
		LEFT JOIN
			phones AS ph ON c.id = ph.company_id AND ph.fax = false
		LEFT JOIN
			phones AS f ON c.id = f.company_id AND f.fax = true
		WHERE
			c.id = $1
		GROUP BY
			c.id
	`, id).Scan(&company.Name, &company.Address, &company.ScopeID, &company.Note, &company.CreatedAt, &company.UpdatedAt, &emails,
		&phones,
		&faxes,
		&practices,
		&contacts)
	if err != nil {
		errmsg("GetCompany select", err)
		return company, err
	}
	practices, err := e.GetPracticeCompany(id)
	if err != nil {
		errmsg("GetPracticeCompany", err)
		return company, err
	}
	company.Practices = practices
	contacts, err := e.GetContactCompany(id)
	if err != nil {
		errmsg("GetContactCompany", err)
		return company, err
	}
	company.Contacts = contacts
	return company, err
}

// CompanyListGet - get all companyes for list
func CompanyListGet() ([]CompanyList, error) {
	var companies []CompanyList
	_, err := pool.Query(&companies, `
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

// CompanySelectGet - get company for contact
func CompanySelectGet(id int64) (SelectItem, error) {
	var company SelectItem
	if id == 0 {
		return company, nil
	}
	err := pool.QueryRow(context.Background(), &Company{}).
		Column("id", "name").
		Where("id = ?", id).
		Select(&company)
	if err != nil {
		errmsg("GetCompanySelect select", err)
	}
	return company, err
}

// CompanySelectGet - get all companyes for select
func CompanySelectGet() ([]SelectItem, error) {
	var companies []SelectItem
	err := pool.QueryRow(context.Background(), &Company{}).
		Column("id", "name").
		Order("name ASC").
		Select(&companies)
	if err != nil {
		errmsg("GetCompanySelectAll select", err)
	}
	return companies, err
}

// CompanyInsert - create new company
func CompanyInsert(company Company) (int64, error) {
	err := pool.Insert(&company)
	if err != nil {
		errmsg("CreateCompany insert", err)
		return 0, err
	}
	_ = e.UpdateCompanyEmails(company)
	_ = e.UpdateCompanyPhones(company, false)
	_ = e.UpdateCompanyPhones(company, true)
	return company.ID, nil
}

// CompanyUpdate - save company changes
func CompanyUpdate(company Company) error {
	err := pool.Update(&company)
	if err != nil {
		errmsg("UpdateCompany update", err)
		return err
	}
	_ = e.UpdateCompanyEmails(company)
	_ = e.UpdateCompanyPhones(company, false)
	_ = e.UpdateCompanyPhones(company, true)
	return nil
}

// CompanyDelete - delete company by id
func CompanyDelete(id int64) error {
	if id == 0 {
		return nil
	}
	err := e.DeleteAllCompanyPhones(id)
	if err != nil {
		errmsg("DeleteCompany DeleteAllCompanyPhones", err)
	}
	_, err = pool.Model(&Company{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteCompany delete", err)
	}
	return err
}

func companyCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			companies (
				id BIGSERIAL PRIMARY KEY,
				name TEXT,
				address TEXT,
				scope_id BIGINT,
				note TEXT,
				created_at timestamp without time zone,
				updated_at
 timestamp without time zone default now(),
				UNIQUE(name, scope_id)
			)
	`
	_, err := pool.Exec(str)
	if err != nil {
		errmsg("companyCreateTable exec", err)
	}
	return err
}
