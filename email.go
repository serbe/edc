package edc

import (
	"context"
	"time"
)

// Email - struct for email
type Email struct {
	ID        int64  `sql:"id"            json:"id"         form:"id"         query:"id"`
	CompanyID int64  `sql:"company_id,pk" json:"company_id" form:"company_id" query:"company_id"`
	ContactID int64  `sql:"contact_id,pk" json:"contact_id" form:"contact_id" query:"contact_id"`
	Email     string `sql:"email"         json:"email"      form:"email"      query:"email"`
	CreatedAt string `sql:"created_at"    json:"-"`
	UpdatedAt string `sql:"updated_at"    json:"-"`
}

// EmailInsert - create new email
func EmailInsert(email Email) (int64, error) {
	email.ID = 0
	err := pool.QueryRow(context.Background(), `
		INSERT INTO emails
		(
			company_id,
			contact_id,
			email,
			created_at,
			updated_at
		)
		VALUES
		(
			$1,
			$2,
			$3,
			$4,
			$5
		)
	`, email.CompanyID, email.ContactID, email.Email, time.Now(), time.Now()).Scan(&email.ID)
	if err != nil {
		errmsg("EmailInsert QueryRow", err)
	}
	return email.ID, nil
}

// EmailCompanyUpdate - update company emails
func EmailCompanyUpdate(id int64, emails []string) error {
	err := EmailCompanyDelete(id)
	if err != nil {
		errmsg("EmailCompanyUpdate DeleteCompanyEmails", err)
		return err
	}
	for i := range emails {
		var email Email
		email.CompanyID = id
		email.Email = emails[i]
		_, err = EmailInsert(email)
		if err != nil {
			errmsg("EmailCompanyUpdate EmailInsert", err)
			return err
		}
	}
	return nil
}

// EmailContactUpdate - update contact emails
func EmailContactUpdate(id int64, emails []string) error {
	err := EmailContactDelete(id)
	if err != nil {
		errmsg("EmailContactUpdate EmailsContactDelete", err)
		return err
	}
	for i := range emails {
		var email Email
		email.ContactID = id
		email.Email = emails[i]
		_, err = EmailInsert(email)
		if err != nil {
			errmsg("EmailContactUpdate EmailInsert", err)
			return err
		}
	}
	return nil
}

// EmailCompanyDelete - delete all emails by company id
func EmailCompanyDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.Exec(context.Background(), `
		DELETE FROM
			emails
		WHERE
			company_id = $1
	`, id)
	if err != nil {
		errmsg("EmailCompanyDelete Exec", err)
	}
	return err
}

// EmailContactDelete - delete all emails by contact id
func EmailContactDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.Exec(context.Background(), `
		DELETE FROM
			emails
		WHERE
			contact_id = $1
	`, id)
	if err != nil {
		errmsg("EmailContactDelete Exec", err)
	}
	return err
}

func emailCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			emails (
				id bigserial primary key,
				company_id bigint,
				contact_id bigint,
				email text,
				created_at timestamp without time zone,
				updated_at timestamp without time zone default now()
			)
	`
	_, err := pool.Exec(context.Background(), str)
	if err != nil {
		errmsg("emailCreateTable Exec", err)
	}
	return err
}
