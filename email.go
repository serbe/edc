package edc

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
	err := pool.Insert(&email)
	if err != nil {
		errmsg("CreateEmail insert", err)
	}
	return email.ID, nil
}

// CompanyEmailsUpdate - update company emails
func CompanyEmailsUpdate(company Company) error {
	err := e.DeleteCompanyEmails(company.ID)
	if err != nil {
		errmsg("UpdateCompanyEmails DeleteCompanyEmails", err)
		return err
	}
	for i := range company.Emails {
		var email Email
		email.CompanyID = company.ID
		email.Email = company.Emails[i]
		_, err = e.CreateEmail(email)
		if err != nil {
			errmsg("UpdateCompanyEmails CreateEmail", err)
			return err
		}
	}
	return nil
}

// ContactEmailsUpdate - update contact emails
func ContactEmailsUpdate(contact Contact) error {
	err := e.DeleteContactEmails(contact.ID)
	if err != nil {
		errmsg("UpdateContactEmails DeleteContactEmails", err)
		return err
	}
	for i := range contact.Emails {
		var email Email
		email.ContactID = contact.ID
		email.Email = contact.Emails[i]
		_, err = e.CreateEmail(email)
		if err != nil {
			errmsg("UpdateContactEmails CreateEmail", err)
			return err
		}
	}
	return nil
}

// EmailUpdate - save email changes
func EmailUpdate(email Email) error {
	err := pool.Update(&email)
	if err != nil {
		errmsg("UpdateEmail update", err)
	}
	return err
}

// EmailDelete - delete email by id
func EmailDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.QueryRow(context.Background(), &Email{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteEmail delete", err)
	}
	return err
}

// CompanyEmailsDelete - delete all emails by company id
func CompanyEmailsDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.QueryRow(context.Background(), &Email{}).
		Where("company_id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteCompanyEmails delete", err)
	}
	return err
}

// ContactEmailsDelete - delete all emails by contact id
func ContactEmailsDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.QueryRow(context.Background(), &Email{}).
		Where("contact_id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteContactEmails delete", err)
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
				updated_at
 timestamp without time zone default now()
			)
	`
	_, err := pool.Exec(str)
	if err != nil {
		errmsg("emailCreateTable exec", err)
	}
	return err
}
