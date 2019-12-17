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

// // GetEmail - get one email by id
// func GetEmail(id int64) (Email, error) {
// 	var email Email
// 	if id == 0 {
// 		return email, nil
// 	}
// 	err := pool.Model(&email).
// 		Where("id = ?", id).
// 		Select()
// 	if err != nil {
// 		errmsg("GetEmail select", err)
// 	}
// 	return email, nil
// }

// // GetEmails - get all emails for list
// func GetEmails() ([]Email, error) {
// 	var emails []Email
// 	err := pool.Model(&emails).
// 		Column("id", "email").
// 		Order("email ASC").
// 		Select()
// 	if err != nil {
// 		errmsg("GetEmailList select", err)
// 	}
// 	return emails, err
// }

// // GetCompanyEmails - get all emails by company id
// func GetCompanyEmails(id int64) ([]Email, error) {
// 	var emails []Email
// 	if id == 0 {
// 		return emails, nil
// 	}
// 	err := pool.Model(&emails).
// 		Column("id", "email").
// 		Order("email ASC").
// 		Where("company_id = ?", id).
// 		Select()
// 	if err != nil {
// 		errmsg("GetCompanyEmails select", err)
// 	}
// 	return emails, err
// }

// // GetContactEmails - get all emails by contact id
// func GetContactEmails(id int64) ([]Email, error) {
// 	var emails []Email
// 	if id == 0 {
// 		return emails, nil
// 	}
// 	err := pool.Model(&emails).
// 		Column("id", "email").
// 		Order("email ASC").
// 		Where("contact_id = ?", id).
// 		Select()
// 	if err != nil {
// 		errmsg("GetContactEmails select", err)
// 	}
// 	return emails, err
// }

// CreateEmail - create new email
func CreateEmail(email Email) (int64, error) {
	email.ID = 0
	err := pool.Insert(&email)
	if err != nil {
		errmsg("CreateEmail insert", err)
	}
	return email.ID, nil
}

// UpdateCompanyEmails - update company emails
func UpdateCompanyEmails(company Company) error {
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

// UpdateContactEmails - update contact emails
func UpdateContactEmails(contact Contact) error {
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

// UpdateEmail - save email changes
func UpdateEmail(email Email) error {
	err := pool.Update(&email)
	if err != nil {
		errmsg("UpdateEmail update", err)
	}
	return err
}

// DeleteEmail - delete email by id
func DeleteEmail(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.Model(&Email{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteEmail delete", err)
	}
	return err
}

// DeleteCompanyEmails - delete all emails by company id
func DeleteCompanyEmails(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.Model(&Email{}).
		Where("company_id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteCompanyEmails delete", err)
	}
	return err
}

// DeleteContactEmails - delete all emails by contact id
func DeleteContactEmails(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.Model(&Email{}).
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
				updated_at timestamp without time zone default now()
			)
	`
	_, err := pool.Exec(str)
	if err != nil {
		errmsg("emailCreateTable exec", err)
	}
	return err
}
