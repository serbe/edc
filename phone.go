package edc

import (
	"context"
	"time"
)

// Phone - struct for phone
type Phone struct {
	ID        int64  `sql:"id"            json:"id"           form:"id"         query:"id"`
	CompanyID int64  `sql:"company_id,pk" json:"company_id"   form:"company_id" query:"company_id"`
	ContactID int64  `sql:"contact_id,pk" json:"contact_id"   form:"contact_id" query:"contact_id"`
	Phone     int64  `sql:"phone"         json:"phone,string" form:"phone"      query:"phone"`
	Fax       bool   `sql:"fax"           json:"fax"          form:"fax"        query:"fax"`
	CreatedAt string `sql:"created_at"    json:"-"`
	UpdatedAt string `sql:"updated_at"    json:"-"`
}

// PhoneInsert - create new phone
func PhoneInsert(phone Phone) (int64, error) {
	phone.ID = 0
	err := pool.QueryRow(context.Background(), `
		INSERT INTO phones
		(
			company_id,
			contact_id,
			phone,
			fax,
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
	`, phone.CompanyID, phone.ContactID, phone.Phone, phone.Fax, time.Now(), time.Now()).Scan(&phone.ID)
	if err != nil {
		errmsg("PhoneInsert QueryRow", err)
	}
	return phone.ID, nil
}

// PhoneCompanyUpdate - update company phones
func PhoneCompanyUpdate(id int64, phones []int64, fax bool) error {
	err := PhoneCompanyDelete(id, fax)
	if err != nil {
		errmsg("PhoneCompanyUpdate PhonesCompanyDelete", err)
		return err
	}
	for i := range phones {
		var phone Phone
		phone.CompanyID = id
		phone.Phone = phones[i]
		phone.Fax = fax
		_, err = PhoneInsert(phone)
		if err != nil {
			errmsg("PhoneCompanyUpdate PhoneInsert", err)
			return err
		}
	}
	return nil
}

// PhoneContactUpdate - update contact phones
func PhoneContactUpdate(id int64, phones []int64, fax bool) error {
	err := PhoneContactDelete(id, fax)
	if err != nil {
		errmsg("PhoneContactUpdate PhonesContactDelete", err)
		return err
	}
	for i := range phones {
		var phone Phone
		phone.ContactID = id
		phone.Phone = phones[i]
		phone.Fax = fax
		_, err = PhoneInsert(phone)
		if err != nil {
			errmsg("PhoneContactUpdate PhoneInsert", err)
			return err
		}
	}
	return nil
}

// PhoneCompanyDelete - delete all unnecessary phones by company id
func PhoneCompanyDelete(id int64, fax bool) error {
	if id == 0 {
		return nil
	}
	_, err := pool.Exec(context.Background(), `
		DELETE FROM
			phones
		WHERE
			company_id = $1
		AND
			fax = $2
	`, id, fax)
	if err != nil {
		errmsg("PhoneCompanyDelete Exec", err)
	}
	return err
}

// PhoneContactDelete - delete all unnecessary phones by contact id
func PhoneContactDelete(id int64, fax bool) error {
	if id == 0 {
		return nil
	}
	_, err := pool.Exec(context.Background(), `
		DELETE FROM
			phones
		WHERE
			contact_id = $1
		AND
			fax = $2
	`, id, fax)
	if err != nil {
		errmsg("PhoneContactDelete Exec", err)
	}
	return err
}

func phoneCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			phones (
				id bigserial primary key,
				contact_id bigint,
				company_id bigint,
				phone bigint,
				fax bool NOT NULL DEFAULT false,
				created_at TIMESTAMP without time zone,
				updated_at TIMESTAMP without time zone default now()
			)
	`
	_, err := pool.Exec(context.Background(), str)
	if err != nil {
		errmsg("phoneCreateTable Exec", err)
	}
	return err
}
