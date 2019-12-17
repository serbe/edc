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

// // EmailGet - get one email by id
// func EmailGet(id int64) (Email, error) {
// 	var email Email
// 	if id == 0 {
// 		return email, nil
// 	}
// 	err := pool.QueryRow(context.Background(), &email).
// 		Where("id = ?", id).
// 		Select()
// 	if err != nil {
// 		errmsg("GetEmail select", err)
// 	}
// 	return email, nil
// }

// // EmailsGet - get all emails for list
// func EmailsGet() ([]Email, error) {
// 	var emails []Email
// 	err := pool.QueryRow(context.Background(), &emails).
// 		Column("id", "email").
// 		Order("email ASC").
// 		Select()
// 	if err != nil {
// 		errmsg("GetEmailList select", err)
// 	}
// 	return emails, err
// }

// // CompanyEmailsGet - get all emails by company id
// func CompanyEmailsGet(id int64) ([]Email, error) {
// 	var emails []Email
// 	if id == 0 {
// 		return emails, nil
// 	}
// 	err := pool.QueryRow(context.Background(), &emails).
// 		Column("id", "email").
// 		Order("email ASC").
// 		Where("company_id = ?", id).
// 		Select()
// 	if err != nil {
// 		errmsg("GetCompanyEmails select", err)
// 	}
// 	return emails, err
// }

// // ContactEmailsGet - get all emails by contact id
// func ContactEmailsGet(id int64) ([]Email, error) {
// 	var emails []Email
// 	if id == 0 {
// 		return emails, nil
// 	}
// 	err := pool.QueryRow(context.Background(), &emails).
// 		Column("id", "email").
// 		Order("email ASC").
// 		Where("contact_id = ?", id).
// 		Select()
// 	if err != nil {
// 		errmsg("GetContactEmails select", err)
// 	}
// 	return emails, err
// }

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
		errmsg("CreateEmail insert", err)
	}
	return email.ID, nil
}

// EmailsCompanyUpdate - update company emails
func EmailsCompanyUpdate(id int64, emails []string) error {
	err := EmailsCompanyDelete(id)
	if err != nil {
		errmsg("EmailsCompanyUpdate DeleteCompanyEmails", err)
		return err
	}
	for i := range emails {
		var email Email
		email.CompanyID = id
		email.Email = emails[i]
		_, err = EmailInsert(email)
		if err != nil {
			errmsg("EmailsCompanyUpdate CreateEmail", err)
			return err
		}
	}
	return nil
}

// ContactEmailsUpdate - update contact emails
func ContactEmailsUpdate(id int64, emails []string) error {
	err := EmailsContactDelete(id)
	if err != nil {
		errmsg("UpdateContactEmails DeleteContactEmails", err)
		return err
	}
	for i := range emails {
		var email Email
		email.ContactID = id
		email.Email = emails[i]
		_, err = EmailInsert(email)
		if err != nil {
			errmsg("UpdateContactEmails CreateEmail", err)
			return err
		}
	}
	return nil
}

// EmailsCompanyDelete - delete all emails by company id
func EmailsCompanyDelete(id int64) error {
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
		errmsg("EmailsCompanyDelete delete", err)
	}
	return err
}

// EmailsContactDelete - delete all emails by contact id
func EmailsContactDelete(id int64) error {
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
		errmsg("EmailsContactDelete delete", err)
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
		errmsg("emailCreateTable exec", err)
	}
	return err
}
