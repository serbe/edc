package edc

import (
	"fmt"
	"log"
	"time"
)

// Practice - struct for practice
type Practice struct {
	ID             int64     `sql:"id" json:"id"`
	Company        Company   `sql:"-"`
	CompanyID      int64     `sql:"company_id, null" json:"company-id"`
	Kind           Kind      `sql:"-"`
	KindID         int64     `sql:"kind_id, null" json:"kind-id"`
	Topic          string    `sql:"topic, null" json:"topic"`
	DateOfPractice time.Time `sql:"date_of_practice, null" json:"date-of-practice"`
	DateStr        string    `sql:"-" json:"date-str"`
	Notes          string    `sql:"notes, null" json:"notes"`
}

// GetPractice - get one practice by id
func (e *EDc) GetPractice(id int64) (practice Practice, err error) {
	if id == 0 {
		return
	}
	_, err = e.db.QueryOne(&practice, "SELECT * FROM practices WHERE id = ? LIMIT 1", id)
	if err != nil {
		log.Println("GetPractice e.db.QueryRow: ", err)
		return practice, fmt.Errorf("GetPractice e.db.QueryRow: %s", err)
	}
	return
}

// GetPracticeAll - get all practices
func (e *EDc) GetPracticeAll() (practices []Practice, err error) {
	_, err = e.db.Query(&practices, "SELECT * FROM practices")
	if err != nil {
		log.Println("GetPracticeAll e.db.Query: ", err)
		return practices, fmt.Errorf("GetPracticeAll e.db.Query: %s", err)
	}
	for i := range practices {
		practices[i].Company, _ = e.GetCompany(practices[i].CompanyID)
		practices[i].Kind, _ = e.GetKind(practices[i].KindID)
		practices[i].DateStr = setStrMonth(practices[i].DateOfPractice)
	}
	return
}

// GetCompanyPractices - get all practices by company id
func (e *EDc) GetCompanyPractices(id int64) (practices []Practice, err error) {
	if id == 0 {
		return
	}
	_, err = e.db.Query(&practices, "SELECT * FROM practices WHERE company_id = ?", id)
	if err != nil {
		log.Println("GetCompanyPractices e.db.Query: ", err)
		return practices, fmt.Errorf("GetCompanyPractices e.db.Query: %s", err)
	}
	return
}

// CreatePractice - create new practice
func (e *EDc) CreatePractice(practice Practice) (err error) {
	err = e.db.Create(&practice)
	if err != nil {
		log.Println("CreatePractice e.db.Create: ", err)
		return fmt.Errorf("CreatePractice e.db.Create: %s", err)
	}
	return
}

// UpdatePractice - save practice changes
func (e *EDc) UpdatePractice(practice Practice) (err error) {
	err = e.db.Update(&practice)
	if err != nil {
		log.Println("UpdatePractice e.db.Update: ", err)
		return fmt.Errorf("UpdatePractice e.db.Update: %s", err)
	}
	return
}

// DeletePractice - delete practice by id
func (e *EDc) DeletePractice(id int64) (err error) {
	if id == 0 {
		return
	}
	_, err = e.db.Exec("DELETE FROM practices WHERE id=?", id)
	if err != nil {
		log.Println("DeletePractice e.db.Exec: ", err)
		return fmt.Errorf("DeletePractice e.db.Exec: %s", err)
	}
	return
}

func (e *EDc) practiceCreateTable() (err error) {
	str := `CREATE TABLE IF NOT EXISTS practices (id bigserial primary key, company_id bigint, kind_id bigint, topic text, date_of_practice date, notes text)`
	_, err = e.db.Exec(str)
	if err != nil {
		log.Println("practiceCreateTable e.db.Exec: ", err)
		return fmt.Errorf("practiceCreateTable e.db.Exec: %s", err)
	}
	return
}
