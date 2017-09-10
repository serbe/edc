package edc

// Certificate - struct for certificate
type Certificate struct {
	ID        int64      `sql:"id"         json:"id"         form:"id"         query:"id"`
	Num       string     `sql:"num"        json:"num"        form:"num"        query:"num"`
	ContactID int64      `sql:"contact_id" json:"contact_id" form:"contact_id" query:"contact_id"`
	Contact   SelectItem `sql:"-"          json:"contact"    form:"contact"    query:"contact"`
	CertDate  string     `sql:"cert_date"  json:"cert_date"  form:"cert_date"  query:"cert_date"`
	Note      string     `sql:"note,null"  json:"note"       form:"note"       query:"note"`
	CreatedAt string     `sql:"created_at" json:"-"`
	UpdatedAt string     `sql:"updated_at" json:"-"`
}

// CertificateList - struct for certificate list
type CertificateList struct {
	ID          int64  `sql:"id"           json:"id"           form:"id"           query:"id"`
	Num         string `sql:"num"          json:"num"          form:"num"          query:"num"`
	ContactID   int64  `sql:"contact_id"   json:"contact_id"   form:"contact_id"   query:"contact_id"`
	ContactName string `sql:"contact_name" json:"contact_name" form:"contact_name" query:"contact_name"`
	CertDate    string `sql:"cert_date"    json:"cert_date"    form:"cert_date"    query:"cert_date"`
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
			co.name AS contact_name,
			c.cert_date
		FROM
			certificates AS c
		LEFT JOIN
			contacts AS co ON c.contact_id = co.id
		GROUP BY
			c.id,
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
