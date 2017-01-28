package edc

// Phone - struct for phone
type Phone struct {
	ID        int64  `sql:"id" json:"id"`
	CompanyID int64  `sql:"company_id, pk, null" json:"company_id"`
	ContactID int64  `sql:"contact_id, pk, null" json:"contact_id"`
	Phone     int64  `sql:"phone, null" json:"phone"`
	Fax       bool   `sql:"fax, null" json:"fax"`
	CreatedAt string `sql:"created_at" json:"created_at"`
	UpdatedAt string `sql:"updated_at" json:"updated_at"`
}

// PhoneSelect - struct for short phone
type PhoneSelect struct {
	ID    int64 `json:"id"`
	Phone int64 `json:"phone"`
}

// GetPhone - get one phone by id
func (e *Edb) GetPhone(id int64) (Phone, error) {
	var phone Phone
	if id == 0 {
		return phone, nil
	}
	err := e.db.Model(&phone).Where("id = ?", id).Select()
	if err != nil {
		errmsg("GetPhone select", err)
	}
	return phone, nil
}

// GetPhoneList - get all phones for list
func (e *Edb) GetPhoneList() ([]Phone, error) {
	var phones []Phone
	_, err := e.db.Query(&phones, `
		SELECT
			id,
			company_id,
			contact_id,
			phone,
			fax
		FROM
			phones
		ORDER BY
			phone ASC`)
	if err != nil {
		errmsg("GetPhoneList query", err)
	}
	return phones, err
}

// GetCompanyPhones - get all phones by company id
func (e *Edb) GetCompanyPhones(id int64, fax bool) ([]PhoneSelect, error) {
	var phones []PhoneSelect
	if id == 0 {
		return phones, nil
	}
	_, err := e.db.Query(&phones, `
		SELECT
			id,
			phone
		FROM
			phones
		WHERE
			company_id = ? AND fax = ?
		ORDER BY
			phone ASC
	`, id, fax)
	if err != nil {
		errmsg("GetCompanyPhones query", err)
	}
	return phones, err
}

// GetContactPhones - get all phones by contact id
func (e *Edb) GetContactPhones(id int64, fax bool) ([]PhoneSelect, error) {
	var phones []PhoneSelect
	if id == 0 {
		return phones, nil
	}
	_, err := e.db.Query(&phones, `
		SELECT
			id,
			phone
		FROM
			phones
		ORDER BY
			phone ASC
		WHERE
			contact_id = ? AND fax = ?
	`, id, fax)
	if err != nil {
		errmsg("GetContactPhones equery", err)
	}
	return phones, nil
}

// GetCompanyPhonesAll - get all faxes or phones by company id and isfax
func (e *Edb) GetCompanyPhonesAll(id int64, fax bool) ([]PhoneSelect, error) {
	var phones []PhoneSelect
	if id == 0 {
		return phones, nil
	}
	_, err := e.db.Query(&phones, `
		SELECT
			id,
			phone
		FROM
			phones
		WHERE
			company_id = ? and fax = ?
		ORDER BY
			phone ASC
	`, id, fax)
	if err != nil {
		errmsg("GetCompanyPhonesAll query", err)
	}
	return phones, err
}

// GetContactPhonesAll - get all phones and faxes by contact id
func (e *Edb) GetContactPhonesAll(id int64, fax bool) ([]PhoneSelect, error) {
	var phones []PhoneSelect
	if id == 0 {
		return phones, nil
	}
	_, err := e.db.Query(&phones, `
		SELECT
			id,
			phone
		FROM
			phones
		WHERE
			contact_id = ? and fax = ?
		ORDER BY
			phone ASC
	`, id, fax)
	if err != nil {
		errmsg("GetContactPhonesAll query", err)
	}
	return phones, err
}

// CreatePhone - create new phone
func (e *Edb) CreatePhone(phone Phone) (int64, error) {
	err := e.db.Insert(&phone)
	if err != nil {
		errmsg("CreatePhone insert", err)
	}
	return phone.ID, nil
}

// CreateCompanyPhones - create new phones to company
func (e *Edb) CreateCompanyPhones(company Company, fax bool) error {
	err := e.CleanCompanyPhones(company, fax)
	if err != nil {
		errmsg("CreateCompanyPhones CleanCompanyPhones", err)
		return err
	}
	for _, value := range company.Phones {
		var id int64
		_, err = e.db.QueryOne(&id, `
			SELECT
				id
			FROM
				phones
			WHERE
				company_id = ? and phone = ? and fax = ?
			RETURNING
				id
		`, company.ID, value.Phone, fax)
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
	return nil
}

// CreateContactPhones - create new phones to contact
func (e *Edb) CreateContactPhones(contact Contact, fax bool) error {
	err := e.CleanContactPhones(contact, fax)
	if err != nil {
		errmsg("CreateContactPhones CleanContactPhones", err)
		return err
	}
	var allPhones []Phone
	if fax {
		allPhones = contact.Faxes
	} else {
		allPhones = contact.Phones
	}
	for _, value := range allPhones {
		var id int64
		_, err = e.db.QueryOne(&id, `
			SELECT
				id
			FROM
				phones
			WHERE
				contact_id = ? and phone = ? and fax = ?
			RETURNING
				id
		`, contact.ID, value.Phone, fax)
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
	return nil
}

// CleanCompanyPhones - delete all unnecessary phones by company id
func (e *Edb) CleanCompanyPhones(company Company, fax bool) error {
	var (
		phones    []int64
		allPhones []Phone
	)
	if fax {
		allPhones = company.Faxes
	} else {
		allPhones = company.Phones
	}
	for _, value := range allPhones {
		phones = append(phones, value.Phone)
	}
	if len(phones) == 0 {
		_, err := e.db.Model(&Phone{}).Where("company_id = ? and fax = ?", company.ID, fax).Delete()
		if err != nil {
			errmsg("CleanCompanyPhones delete", err)
			return err
		}
	} else {
		var companyPhones []Phone
		_, err := e.db.Query(&companyPhones, `
			SELECT
				id,
				phone
			FROM
				phones
			WHERE
				company_id = ? and fax = ?
		`, company.ID, fax)
		if err != nil {
			errmsg("CleanCompanyPhones query", err)
			return err
		}
		for _, value := range companyPhones {
			if int64InSlice(value.Phone, phones) == false {
				_, err = e.db.Model(&Phone{}).
					Where("company_id = ? and phone = ? and fax = ?", company.ID, value.Phone, fax).
					Delete()
				if err != nil {
					errmsg("CleanCompanyPhones delete", err)
					return err
				}
			}
		}
	}
	return nil
}

// CleanContactPhones - delete all unnecessary phones by contact id
func (e *Edb) CleanContactPhones(contact Contact, fax bool) error {
	var (
		phones    []int64
		allPhones []Phone
	)
	if fax {
		allPhones = contact.Faxes
	} else {
		allPhones = contact.Phones
	}
	for _, value := range allPhones {
		phones = append(phones, value.Phone)
	}
	if len(phones) == 0 {
		_, err := e.db.Model(&Phone{}).Where("contact_id = ? and fax = ?", contact.ID, fax).Delete()
		if err != nil {
			errmsg("CleanContactPhones delete", err)
			return err
		}
	} else {
		var contactPhones []Phone
		_, err := e.db.Query(&contactPhones, `
			SELECT
				id,
				phone
			FROM
				phones
			WHERE
				contact_id = ? and fax = ?
		`, contact.ID, fax)
		if err != nil {
			errmsg("CleanContactPhones query", err)
			return err
		}
		for _, value := range contactPhones {
			if int64InSlice(value.Phone, phones) == false {
				_, err = e.db.Model(&Phone{}).
					Where("contact_id = ? and phone = ? and fax = ?", contact.ID, value.Phone, fax).
					Delete()
				if err != nil {
					errmsg("CleanContactPhones delete", err)
					return err
				}
			}
		}
	}
	return nil
}

// DeleteAllCompanyPhones - delete all phones and faxes by company id
func (e *Edb) DeleteAllCompanyPhones(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Model(&Phone{}).Where("company_id = ?", id).Delete()
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
	_, err := e.db.Model(&Phone{}).Where("contact_id = ?", id).Delete()
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
