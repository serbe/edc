package edc

import (
	"context"
	"time"
)

// Certificate - struct for certificate
type Certificate struct {
	ID        int64  `sql:"id"         json:"id"         form:"id"         query:"id"`
	Num       string `sql:"num"        json:"num"        form:"num"        query:"num"`
	ContactID int64  `sql:"contact_id" json:"contact_id" form:"contact_id" query:"contact_id"`
	CompanyID int64  `sql:"company_id" json:"company_id" form:"company_id" query:"company_id"`
	CertDate  string `sql:"cert_date"  json:"cert_date"  form:"cert_date"  query:"cert_date"`
	Note      string `sql:"note"       json:"note"       form:"note"       query:"note"`
	CreatedAt string `sql:"created_at" json:"-"`
	UpdatedAt string `sql:"updated_at" json:"-"`
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

// CertificateGet - get one certificate by id
func CertificateGet(id int64) (Certificate, error) {
	var certificate Certificate
	if id == 0 {
		return certificate, nil
	}
	certificate.ID = id
	err := pool.QueryRow(context.Background(), `
		SELECT
			num,
			contact_id,
			company_id,
			cert_date,
			note,
			created_at,
			updated_at
 		FROM
			certificates
		WHERE
			id = $1
	`, id).Scan(&certificate.Num, &certificate.ContactID, &certificate.CompanyID, &certificate.CertDate, &certificate.Note, &certificate.CreatedAt, &certificate.UpdatedAt)
	if err != nil {
		errmsg("CertificateGet select", err)
	}
	return certificate, err
}

// CertificateListGet - get all certificate for list
func CertificateListGet() ([]CertificateList, error) {
	var certificates []CertificateList
	rows, err := pool.Query(context.Background(), `
		SELECT
			c.id,
			c.num,
			c.contact_id,
			p.name AS contact_name,
			c.company_id,
			co.name AS company_name,
			c.cert_date,
			c.note
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
		errmsg("CertificateListGet select", err)
	}
	for rows.Next() {
		var certificate CertificateList
		err := rows.Scan(&certificate.ID, &certificate.Num, &certificate.ContactID, &certificate.ContactName, &certificate.CompanyID, &certificate.CompanyName, &certificate.CertDate, &certificate.Note)
		if err != nil {
			errmsg("CertificateListGet select", err)
			return certificates, err
		}
		certificates = append(certificates, certificate)
	}
	return certificates, rows.Err()
}

// CertificateCreate - create new certificate
func CertificateCreate(certificate Certificate) (int64, error) {
	err := pool.QueryRow(context.Background(), `
		INSERT INTO certificates
		(
			num,
			contact_id,
			company_id,
			cert_date,
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
			$6,
			$7
		)
		RETURNING
			id
	`, certificate.Num,
		certificate.ContactID,
		certificate.CompanyID,
		certificate.CertDate,
		certificate.Note,
		time.Now(),
		time.Now()).Scan(&certificate.ID)
	if err != nil {
		errmsg("CertificateCreate insert", err)
	}
	return certificate.ID, nil
}

// CertificateUpdate - save certificate changes
func CertificateUpdate(certificate Certificate) error {
	_, err := pool.Exec(context.Background(), `
		UPDATE certificates SET
			num = $2,
			contact_id = $3,
			company_id = $4,
			cert_date = $5,
			note = $6,
			updated_at = $7
		WHERE
			id = $1
	`, certificate.ID, certificate.Num,
		certificate.ContactID,
		certificate.CompanyID,
		certificate.CertDate,
		certificate.Note,
		time.Now())
	if err != nil {
		errmsg("CertificateUpdate update", err)
	}
	return err
}

// CertificateDelete - delete certificate by id
func CertificateDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.Exec(context.Background(), `
		DELETE FROM
			certificates
		WHERE
			id = $1
	`, id)
	if err != nil {
		errmsg("CertificateDelete delete", err)
	}
	return err
}

func certificateCreateTable() error {
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
	_, err := pool.Exec(context.Background(), str)
	if err != nil {
		errmsg("certificateCreateTable exec", err)
	}
	return err
}
