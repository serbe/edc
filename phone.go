package edc

import "log"

// Phone - struct for phone
type Phone struct {
	ID        int64  `sql:"id" json:"id"`
	PeopleID  int64  `sql:"people_id, pk, null" json:"people-id"`
	CompanyID int64  `sql:"company_id, pk, null" json:"company-id"`
	Phone     int64  `sql:"phone, null" json:"phone"`
	Fax       bool   `sql:"fax, null" json:"fax"`
	Notes     string `sql:"notes, null" json:"notes"`
}

// CreatePhone - create new phone
func (e *EDc) CreatePhone(phone Phone, fax bool) (err error) {
	phone.Fax = fax
	err = e.db.Create(&phone)
	if err != nil {
		log.Println("CreatePhone e.db.Exec ", err)
	}
	return
}

// GetCompanyPhones - get all phones by company id
func (e *EDc) GetCompanyPhones(id int64) (phones []Phone, err error) {
	if id == 0 {
		return
	}
	phones, err = e.GetCompanyPhonesAndFaxes(id, false)
	if err != nil {
		log.Println("GetCompanyPhones GetCompanyPhonesAndFaxes ", err)
	}
	return
}

// GetCompanyFaxes - get all faxes by company id
func (e *EDc) GetCompanyFaxes(id int64) (phones []Phone, err error) {
	if id == 0 {
		return
	}
	phones, err = e.GetCompanyPhonesAndFaxes(id, true)
	if err != nil {
		log.Println("GetCompanyFaxes GetCompanyPhonesAndFaxes ", err)
	}
	return
}

// GetCompanyPhonesAndFaxes - get all faxes or phones by company id and isfax
func (e *EDc) GetCompanyPhonesAndFaxes(id int64, fax bool) (phones []Phone, err error) {
	if id == 0 {
		return
	}
	_, err = e.db.Query(&phones, "SELECT * FROM phones WHERE company_id = ? and fax = ?", id, fax)
	if err != nil {
		log.Println("GetCompanyPhonesAndFaxes e.db.Query ", err)
		return
	}
	return
}

// CreateCompanyPhones - create new phones to company
func (e *EDc) CreateCompanyPhones(company Company) (err error) {
	err = e.CleanCompanyPhones(company)
	if err != nil {
		log.Println("CreateCompanyPhones CleanCompanyPhones ", err)
		return
	}
	for _, value := range company.Phones {
		phone := Phone{}
		_, err = e.db.QueryOne(&phone, "SELECT * FROM phones WHERE company_id = ? and phone = ? and fax = ? LIMIT 1", company.ID, value.Phone, false)
		if phone.ID == 0 {
			value.CompanyID = company.ID
			err = e.CreatePhone(value, false)
			if err != nil {
				log.Println("CreateCompanyPhones CreatePhone ", err)
				return
			}
		}
	}
	return
}

// CreateCompanyFaxes - create new faxes to company
func (e *EDc) CreateCompanyFaxes(company Company) (err error) {
	err = e.CleanCompanyFaxes(company)
	if err != nil {
		log.Println("CreateCompanyFaxes CleanCompanyFaxes ", err)
		return
	}
	for _, value := range company.Faxes {
		phone := Phone{}
		_, err = e.db.QueryOne(&company, "SELECT * FROM phones WHERE company_id = ? and phone = ? and fax = ? LIMIT 1", company.ID, value.Phone, true)
		if phone.ID == 0 {
			value.CompanyID = company.ID
			err = e.CreatePhone(value, true)
			if err != nil {
				log.Println("CreateCompanyFaxes CreatePhone ", err)
				return
			}
		}
	}
	return
}

// CleanCompanyPhones - delete all unnecessary phones by company id
func (e *EDc) CleanCompanyPhones(company Company) (err error) {
	phones := []int64{}
	for _, value := range company.Phones {
		phones = append(phones, value.Phone)
	}
	if len(phones) == 0 {
		_, err = e.db.Exec("DELETE FROM phones WHERE company_id = ? and fax = ?", company.ID, false)
		if err != nil {
			log.Println("CleanCompanyPhones e.db.Exec ", err)
			return
		}
	} else {
		var companyPhones []Phone
		_, err = e.db.Query(&companyPhones, "SELECT * FROM phones WHERE company_id = ? and fax = ?", company.ID, false)
		if err != nil {
			log.Println("CleanCompanyPhones e.db.Query ", err)
			return
		}
		for _, value := range companyPhones {
			if int64InSlice(value.Phone, phones) == false {
				_, err = e.db.Exec("DELETE FROM phones WHERE company_id = ? and fax = ? and phone = ?", company.ID, false, value.Phone)
				if err != nil {
					log.Println("CleanCompanyPhones e.db.Exec ", err)
					return
				}
			}
		}
	}
	return
}

// CleanCompanyFaxes - delete all unnecessary faxes by company id
func (e *EDc) CleanCompanyFaxes(company Company) (err error) {
	phones := []int64{}
	for _, value := range company.Faxes {
		phones = append(phones, value.Phone)
	}
	if len(phones) == 0 {
		_, err = e.db.Exec("DELETE FROM phones WHERE company_id = ? and fax = ?", company.ID, true)
		if err != nil {
			log.Println("CleanCompanyFaxes e.db.Exec ", err)
			return
		}
	} else {
		var companyPhones []Phone
		_, err = e.db.Query(&companyPhones, "SELECT * FROM phones WHERE company_id = ? and fax = ?", company.ID, true)
		if err != nil {
			log.Println("CleanCompanyFaxes e.db.Query ", err)
			return
		}
		for _, value := range companyPhones {
			if int64InSlice(value.Phone, phones) == false {
				_, err = e.db.Exec("DELETE FROM phones WHERE company_id = ? and fax = ? and phone = ?", company.ID, true, value.Phone)
				if err != nil {
					log.Println("CleanCompanyFaxes e.db.Exec ", err)
					return
				}
			}
		}
	}
	return
}

// DeleteAllCompanyPhones - delete all phones and faxes by company id
func (e *EDc) DeleteAllCompanyPhones(id int64) (err error) {
	if id == 0 {
		return
	}
	_, err = e.db.Exec("DELETE FROM phones WHERE company_id = ?", id)
	if err != nil {
		log.Println("DeleteAllCompanyPhones e.db.Exec ", err)
	}
	return
}

// GetPeoplePhones - get all phones by people id
func (e *EDc) GetPeoplePhones(id int64) (phones []Phone, err error) {
	if id == 0 {
		return
	}
	phones, err = e.GetPeoplePhonesAndFaxes(id, false)
	if err != nil {
		log.Println("GetPeoplePhones GetPeoplePhonesAndFaxes ", err)
	}
	return
}

// GetPeopleFaxes - get all faxes by people id
func (e *EDc) GetPeopleFaxes(id int64) (phones []Phone, err error) {
	if id == 0 {
		return
	}
	phones, err = e.GetPeoplePhonesAndFaxes(id, true)
	if err != nil {
		log.Println("GetPeopleFaxes GetPeoplePhonesAndFaxes ", err)
	}
	return
}

// GetPeoplePhonesAndFaxes - get all phones and faxes by people id
func (e *EDc) GetPeoplePhonesAndFaxes(id int64, fax bool) (phones []Phone, err error) {
	if id == 0 {
		return
	}
	_, err = e.db.Query(&phones, "SELECT * FROM phones WHERE people_id = ? and fax = ?", id, fax)
	if err != nil {
		log.Println("GetPeoplePhonesAndFaxes e.db.Query ", err)
		return
	}
	return
}

// CreatePeoplePhones - create new phones to people
func (e *EDc) CreatePeoplePhones(people People) (err error) {
	err = e.CleanPeoplePhones(people)
	if err != nil {
		log.Println("CreatePeoplePhones CleanPeoplePhones ", err)
		return
	}
	for _, value := range people.Phones {
		phone := Phone{}
		_, err = e.db.QueryOne(&phone, "SELECT * FROM phones WHERE people_id = ? and phone = ? and fax = ? LIMIT 1", people.ID, value.Phone, false)
		if phone.ID == 0 {
			value.PeopleID = people.ID
			err = e.CreatePhone(value, false)
			if err != nil {
				log.Println("CreatePeoplePhones CreatePhone ", err)
				return
			}
		}
	}
	return
}

// CreatePeopleFaxes - create new faxes to people
func (e *EDc) CreatePeopleFaxes(people People) (err error) {
	err = e.CleanPeopleFaxes(people)
	if err != nil {
		log.Println("CreatePeopleFaxes CleanPeopleFaxes ", err)
		return
	}
	for _, value := range people.Faxes {
		phone := Phone{}
		_, err = e.db.QueryOne(&phone, "SELECT * FROM phones WHERE people_id = ? and phone = ? and fax = ? LIMIT 1", people.ID, value.Phone, true)
		if phone.ID == 0 {
			value.PeopleID = people.ID
			err = e.CreatePhone(value, true)
			if err != nil {
				log.Println("CreatePeopleFaxes CreatePhone ", err)
				return
			}
		}
	}
	return
}

// CleanPeoplePhones - delete all unnecessary phones by people id
func (e *EDc) CleanPeoplePhones(people People) (err error) {
	phones := []int64{}
	for _, value := range people.Phones {
		phones = append(phones, value.Phone)
	}
	if len(phones) == 0 {
		_, err = e.db.Exec("DELETE FROM phones WHERE people_id = ? and fax = ?", people.ID, false)
		if err != nil {
			log.Println("CleanPeoplePhones e.db.Exec ", err)
			return
		}
	} else {
		var peoplePhones []Phone
		_, err = e.db.Query(&peoplePhones, "SELECT * FROM phones WHERE people_id = ? and fax = ?", people.ID, false)
		if err != nil {
			log.Println("CleanPeoplePhones e.db.Query ", err)
			return
		}
		for _, value := range peoplePhones {
			if int64InSlice(value.Phone, phones) == false {
				_, err = e.db.Exec("DELETE FROM phones WHERE people_id = ? and fax = ? and phone = ?", people.ID, false, value.Phone)
				if err != nil {
					log.Println("CleanPeoplePhones e.db.Exec ", err)
					return
				}
			}
		}
	}
	return
}

// CleanPeopleFaxes - delete all unnecessary faxes by people id
func (e *EDc) CleanPeopleFaxes(people People) (err error) {
	phones := []int64{}
	for _, value := range people.Faxes {
		phones = append(phones, value.Phone)
	}
	if len(phones) == 0 {
		_, err = e.db.Exec("DELETE FROM phones WHERE people_id = ? and fax = ?", people.ID, true)
		if err != nil {
			log.Println("CleanPeopleFaxes e.db.Exec ", err)
			return
		}
	} else {
		var peoplePhones []Phone
		_, err = e.db.Query(&peoplePhones, "SELECT * FROM phones WHERE people_id = ? and fax = ?", people.ID, true)
		if err != nil {
			log.Println("CleanPeopleFaxes e.db.Query ", err)
			return
		}
		for _, value := range peoplePhones {
			if int64InSlice(value.Phone, phones) == false {
				_, err = e.db.Exec("DELETE FROM phones WHERE people_id = ? and fax = ? and phone = ?", people.ID, true, value.Phone)
				if err != nil {
					log.Println("CleanPeopleFaxes e.db.Exec ", err)
					return
				}
			}
		}
	}
	return
}

// DeleteAllPeoplePhones - delete all phones and faxes by people id
func (e *EDc) DeleteAllPeoplePhones(id int64) (err error) {
	if id == 0 {
		return
	}
	_, err = e.db.Exec("DELETE FROM phones WHERE people_id = ?", id)
	if err != nil {
		log.Println("DeleteAllPeoplePhones e.db.Exec ", err)
	}
	return
}

func (e *EDc) phoneCreateTable() (err error) {
	str := `CREATE TABLE IF NOT EXISTS phones (id bigserial primary key, people_id bigint, company_id bigint, phone bigint, fax bool NOT NULL DEFAULT false, notes text)`
	_, err = e.db.Exec(str)
	if err != nil {
		log.Println("phoneCreateTable e.db.Exec ", err)
	}
	return
}
