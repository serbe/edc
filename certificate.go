package edc

// Certificate - struct for certificate
type Certificate struct {
	ID        int64      `sql:"id"         json:"id"         form:"id"         query:"id"`
	Num       string     `sql:"num"        json:"num"        form:"num"        query:"num"`
	ContactID int64      `sql:"contact_id" json:"contact_id" form:"contact_id" query:"contact_id"`
	CompanyID int64      `sql:"company_id" json:"company_id" form:"company_id" query:"company_id"`
	CertDate  string     `sql:"cert_date"  json:"cert_date"  form:"cert_date"  query:"cert_date"`
	Note      string     `sql:"note"       json:"note"       form:"note"       query:"note"`
	CreatedAt string     `sql:"created_at" json:"-"`
	UpdatedAt string     `sql:"updated_at" json:"-"`
}

// CertificateList - struct for certificate list
type CertificateList struct {
	ID          int64  `sql:"id"           json:"id"           form:"id"           query:"id"`
	Num         string `sql:"num"          json:"num"          form:"num"          query:"num"`
	ContactID   int64  `sql:"contact_id"   json:"contact_id"   form:"contact_id"   query:"contact_id"`
	ContactName string `sql:"contact_name" json:"contact_name" form:"contact_name" query:"contact_name"`
	CompanyID   int64  `sql:"company_id"   json:"company_id"   form:"company_id"   query:"company_id"`
	CompanyName string `sql:"company_name" json:"company_name" form:"company_name" query:"company_name"`
	CertDate    string `sql:"cert_date"    json:"cert_date"    form:"cert_date"    query:"cert_date"`
	Note        string `sql:"note"         json:"note"         form:"note"         query:"note"`
}

// GetCertificate - get one certificate by id
func (e *Edb) GetCertificate(id int64) (Certificate, error) {
	var certificate Certificate
	if id == 0 {
		return certificate, nil
	}
	err := e.db.Model(&certificate).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetCertificate select", err)
	}
	return certificate, err
}

// GetCertificateListAll - get all certificate for list
func (e *Edb) GetCertificateListAll() ([]CertificateList, error) {
	var certificates []CertificateList
	_, err := e.db.Query(&certificates, `
		SELECT
			c.id,
			c.num,
			c.contact_id,
			p.name AS contact_name,
			c.company_id,
			co.name AS company_name,
			c.cert_date
		FROM
			certificates AS c
		LEFT JOIN
			contacts AS p ON c.contact_id = p.id
		LEFT JOIN
			companies AS co ON c.company_id = co.id
		GROUP BY
			c.id,
			p.name,
			co.name
		ORDER BY
			num ASC
	`)
	if err != nil {
		errmsg("GetCertificateListAll select", err)
	}
	return certificates, err
}

// CreateCertificate - create new certificate
func (e *Edb) CreateCertificate(certificate Certificate) (int64, error) {
	err := e.db.Insert(&certificate)
	if err != nil {
		errmsg("CreateCertificate insert", err)
	}
	return certificate.ID, nil
}

// UpdateCertificate - save certificate changes
func (e *Edb) UpdateCertificate(certificate Certificate) error {
	err := e.db.Update(&certificate)
	if err != nil {
		errmsg("UpdateCertificate update", err)
	}
	return err
}

// DeleteCertificate - delete certificate by id
func (e *Edb) DeleteCertificate(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Model(&Certificate{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteCertificate delete", err)
	}
	return err
}

func (e *Edb) certificateCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			certificates (
				id BIGSERIAL PRIMARY KEY,
				num TEXT,
				contact_id BIGINT,
				company_id BIGINT,
				cert_date DATE,
				note TEXT,
				created_at TIMESTAMP without time zone,
				updated_at TIMESTAMP without time zone default now(),
				UNIQUE(num)
			)
	`
	_, err := e.db.Exec(str)
	if err != nil {
		errmsg("certificateCreateTable exec", err)
	}
	return err
}
