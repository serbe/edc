package edc

// Phone - struct for phone
type Phone struct {
	ID        int64  `sql:"id"                   json:"id"`
	CompanyID int64  `sql:"company_id, pk, null" json:"company_id"`
	ContactID int64  `sql:"contact_id, pk, null" json:"contact_id"`
	Phone     int64  `sql:"phone, null"          json:"phone, string"`
	Fax       bool   `sql:"fax, null"            json:"fax"`
	CreatedAt string `sql:"created_at"           json:"-"`
	UpdatedAt string `sql:"updated_at"           json:"-"`
}

// GetPhone - get one phone by id
func (e *Edb) GetPhone(id int64) (Phone, error) {
	var phone Phone
	if id == 0 {
		return phone, nil
	}
	err := e.db.Model(&phone).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetPhone select", err)
	}
	return phone, nil
}

// GetPhoneList - get all phones for list
func (e *Edb) GetPhoneList() ([]Phone, error) {
	var phones []Phone
	err := e.db.Model(&phones).
		Column("id", "company_id", "contact_id", "phone", "fax").
		Order("phone ASC").
		Select()
	if err != nil {
		errmsg("GetPhoneList select", err)
	}
	return phones, err
}

// GetCompanyPhones - get all phones by company id
func (e *Edb) GetCompanyPhones(id int64, fax bool) ([]Phone, error) {
	var phones []Phone
	if id == 0 {
		return phones, nil
	}
	err := e.db.Model(&phones).
		Where("company_id = ? AND fax = ?", id, fax).
		Order("phone ASC").
		Select()
	if err != nil {
		errmsg("GetCompanyPhones select", err)
	}
	return phones, err
}

// GetContactPhones - get all phones by contact id
func (e *Edb) GetContactPhones(id int64, fax bool) ([]Phone, error) {
	var phones []Phone
	if id == 0 {
		return phones, nil
	}
	err := e.db.Model(&phones).
		Where("contact_id = ? AND fax = ?", id, fax).
		Order("phone ASC").
		Select()
	if err != nil {
		errmsg("GetContactPhones select", err)
	}
	return phones, nil
}

// CreatePhone - create new phone
func (e *Edb) CreatePhone(phone Phone) (int64, error) {
	phone.ID = 0
	err := e.db.Insert(&phone)
	if err != nil {
		errmsg("CreatePhone insert", err)
	}
	return phone.ID, nil
}

// CreateCompanyPhones - create new phones to company
func (e *Edb) CreateCompanyPhones(company Company, fax bool) error {
	err := e.DeleteCompanyPhones(company.ID, fax)
	if err != nil {
		errmsg("CreateCompanyPhones DeleteCompanyPhones", err)
		return err
	}
	var allPhones []Phone
	if fax {
		allPhones = company.Faxes
	} else {
		allPhones = company.Phones
	}
	for _, value := range allPhones {
		if value.Phone != 0 {
			var id int64
			err = e.db.Model(&Phone{}).
				Column("id").
				Where("company_id = ? and phone = ? and fax = ?", company.ID, value.Phone, fax).
				Select(&id)
			if id == 0 {
				value.CompanyID = company.ID
				value.Fax = fax
				_, err = e.CreatePhone(value)
				if err != nil {
					errmsg("CreateCompanyPhones CreatePhone", err)
					return err
				}
			}
		}
	}
	return nil
}

// CreateContactPhones - create new phones to contact
func (e *Edb) CreateContactPhones(contact Contact, fax bool) error {
	err := e.DeleteContactPhones(contact.ID, fax)
	if err != nil {
		errmsg("CreateContactPhones DeleteContactPhones", err)
		return err
	}
	var allPhones []Phone
	if fax {
		allPhones = contact.Faxes
	} else {
		allPhones = contact.Phones
	}
	for _, value := range allPhones {
		if value.Phone != 0 {
			var id int64
			err = e.db.Model(&Phone{}).
				Column("id").
				Where("contact_id = ? and phone = ? and fax = ?", contact.ID, value.Phone, fax).
				Select(&id)
			if id == 0 {
				value.ContactID = contact.ID
				value.Fax = fax
				_, err = e.CreatePhone(value)
				if err != nil {
					errmsg("CreateContactPhones CreatePhone", err)
					return err
				}
			}
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
