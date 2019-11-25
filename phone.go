package edc

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

// // GetPhone - get one phone by id
// func (e *Edb) GetPhone(id int64) (Phone, error) {
// 	var phone Phone
// 	if id == 0 {
// 		return phone, nil
// 	}
// 	err := e.db.Model(&phone).
// 		Where("id = ?", id).
// 		Select()
// 	if err != nil {
// 		errmsg("GetPhone select", err)
// 	}
// 	return phone, nil
// }

// // GetPhoneList - get all phones for list
// func (e *Edb) GetPhoneList() ([]Phone, error) {
// 	var phones []Phone
// 	err := e.db.Model(&phones).
// 		Column("id", "company_id", "contact_id", "phone", "fax").
// 		Order("phone ASC").
// 		Select()
// 	if err != nil {
// 		errmsg("GetPhoneList select", err)
// 	}
// 	return phones, err
// }

// // GetCompanyPhones - get all phones by company id
// func (e *Edb) GetCompanyPhones(id int64, fax bool) ([]Phone, error) {
// 	var phones []Phone
// 	if id == 0 {
// 		return phones, nil
// 	}
// 	err := e.db.Model(&phones).
// 		Where("company_id = ? AND fax = ?", id, fax).
// 		Order("phone ASC").
// 		Select()
// 	if err != nil {
// 		errmsg("GetCompanyPhones select", err)
// 	}
// 	return phones, err
// }

// // GetContactPhones - get all phones by contact id
// func (e *Edb) GetContactPhones(id int64, fax bool) ([]Phone, error) {
// 	var phones []Phone
// 	if id == 0 {
// 		return phones, nil
// 	}
// 	err := e.db.Model(&phones).
// 		Where("contact_id = ? AND fax = ?", id, fax).
// 		Order("phone ASC").
// 		Select()
// 	if err != nil {
// 		errmsg("GetContactPhones select", err)
// 	}
// 	return phones, nil
// }

// CreatePhone - create new phone
func (e *Edb) CreatePhone(phone Phone) (int64, error) {
	phone.ID = 0
	err := e.db.Insert(&phone)
	if err != nil {
		errmsg("CreatePhone insert", err)
	}
	return phone.ID, nil
}

// UpdateCompanyPhones - update company phones
func (e *Edb) UpdateCompanyPhones(company Company, fax bool) error {
	err := e.DeleteCompanyPhones(company.ID, fax)
	if err != nil {
		errmsg("UpdateCompanyPhones DeleteCompanyPhones", err)
		return err
	}
	var allPhones []int64
	if fax {
		allPhones = company.Faxes
	} else {
		allPhones = company.Phones
	}
	for i := range allPhones {
		var phone Phone
		phone.CompanyID = company.ID
		phone.Phone = allPhones[i]
		phone.Fax = fax
		_, err = e.CreatePhone(phone)
		if err != nil {
			errmsg("UpdateCompanyPhones CreatePhone", err)
			return err
		}
	}
	return nil
}

// UpdateContactPhones - update contact phones
func (e *Edb) UpdateContactPhones(contact Contact, fax bool) error {
	err := e.DeleteContactPhones(contact.ID, fax)
	if err != nil {
		errmsg("CreateContactPhones DeleteContactPhones", err)
		return err
	}
	var allPhones []int64
	if fax {
		allPhones = contact.Faxes
	} else {
		allPhones = contact.Phones
	}
	for i := range allPhones {
		var phone Phone
		phone.ContactID = contact.ID
		phone.Phone = allPhones[i]
		phone.Fax = fax
		_, err = e.CreatePhone(phone)
		if err != nil {
			errmsg("UpdateContactPhones CreatePhone", err)
			return err
		}
	}
	return nil
}

// DeleteCompanyPhones - delete all unnecessary phones by company id
func (e *Edb) DeleteCompanyPhones(id int64, fax bool) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Model(&Phone{}).
		Where("company_id = ? and fax = ?", id, fax).
		Delete()
	if err != nil {
		errmsg("DeleteCompanyPhones delete", err)
	}
	return err
}

// DeleteContactPhones - delete all unnecessary phones by contact id
func (e *Edb) DeleteContactPhones(id int64, fax bool) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Model(&Phone{}).
		Where("contact_id = ? and fax = ?", id, fax).
		Delete()
	if err != nil {
		errmsg("DeleteContactPhones delete", err)
	}
	return err
}

// DeleteAllCompanyPhones - delete all phones and faxes by company id
func (e *Edb) DeleteAllCompanyPhones(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Model(&Phone{}).
		Where("company_id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteAllCompanyPhones delete", err)
	}
	return err
}

// DeleteAllContactPhones - delete all phones and faxes by contact id
func (e *Edb) DeleteAllContactPhones(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Model(&Phone{}).
		Where("contact_id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteAllContactPhones delete", err)
	}
	return err
}

func (e *Edb) phoneCreateTable() error {
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
	_, err := e.db.Exec(str)
	if err != nil {
		errmsg("phoneCreateTable exec", err)
	}
	return err
}
