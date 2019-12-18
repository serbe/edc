package edc

import (
	"context"
	"time"
)

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
		errmsg("GetCompany QueryRow", err)
		return company, err
	}
	practices, err = PracticeCompanyGet(id)
	if err != nil {
		errmsg("PracticeCompanyGet", err)
		return company, err
	}
	company.Practices = practices
	contacts, err = ContactCompanyGet(id)
	if err != nil {
		errmsg("ContactCompanyGet", err)
		return company, err
	}
	company.Contacts = contacts
	return company, err
}

// CompanyListGet - get all companyes for list
func CompanyListGet() ([]CompanyList, error) {
	var companies []CompanyList
	rows, err := pool.Query(context.Background(), `
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
		errmsg("GetCompanyList Query", err)
	}
	for rows.Next() {
		var company CompanyList
		err := rows.Scan(&company.ID, &company.Name, &company.Address, &company.ScopeName,
			&company.Emails, &company.Phones, &company.Faxes, &company.Practices)
		if err != nil {
			errmsg("GetCompanyList Scan", err)
			return companies, err
		}
		companies = append(companies, company)
	}
	return companies, rows.Err()
}

// CompanySelectGet - get all companyes for select
func CompanySelectGet() ([]SelectItem, error) {
	var companies []SelectItem
	rows, err := pool.Query(context.Background(), `
		SELECT
			id,
			name
		FROM
			companies
		ORDER BY
			name ASC
	`)
	if err != nil {
		errmsg("CompanySelectGet Query", err)
	}
	for rows.Next() {
		var company SelectItem
		err := rows.Scan(&company.ID, &company.Name)
		if err != nil {
			errmsg("CompanySelectGet Scan", err)
			return companies, err
		}
		companies = append(companies, company)
	}
	return companies, rows.Err()
}

// CompanyInsert - create new company
func CompanyInsert(company Company) (int64, error) {
	err := pool.QueryRow(context.Background(), `
		INSERT INTO companies
		(
			name,
			address,
			scope_id,
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
			$6
		)
		RETURNING
			id
	`, company.Name,
		company.Address,
		company.ScopeID,
		company.Note,
		time.Now(),
		time.Now()).Scan(&company.ID)
	if err != nil {
		errmsg("CreateCompany QueryRow", err)
		return 0, err
	}
	_ = EmailsCompanyUpdate(company.ID, company.Emails)
	_ = PhonesCompanyUpdate(company.ID, company.Phones, false)
	_ = PhonesCompanyUpdate(company.ID, company.Faxes, true)
	return company.ID, nil
}

// CompanyUpdate - save company changes
func CompanyUpdate(company Company) error {
	_, err := pool.Exec(context.Background(), `
		UPDATE companies SET
			name = $2,
			address = $3,
			scope_id = $4,
			note = $5,
			updated_at = $6
		WHERE
			id = $1
	`, company.ID, company.Name,
		company.Address,
		company.ScopeID,
		company.Note,
		time.Now())
	if err != nil {
		errmsg("CompanyUpdate Exec", err)
		return err
	}
	_ = EmailsCompanyUpdate(company.ID, company.Emails)
	_ = PhonesCompanyUpdate(company.ID, company.Phones, false)
	_ = PhonesCompanyUpdate(company.ID, company.Faxes, true)
	return nil
}

// CompanyDelete - delete company by id
func CompanyDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.Exec(context.Background(), `
		DELETE FROM
			companyes
		WHERE
			id = $1
	`, id)
	if err != nil {
		errmsg("DeleteCompany Exec", err)
	}
	_ = EmailsCompanyDelete(id)
	_ = PhonesCompanyDelete(id, false)
	_ = PhonesCompanyDelete(id, true)
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
				updated_at timestamp without time zone default now(),
				UNIQUE(name, scope_id)
			)
	`
	_, err := pool.Exec(context.Background(), str)
	if err != nil {
		errmsg("companyCreateTable Exec", err)
	}
	return err
}
