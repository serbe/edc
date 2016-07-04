package edc

import (
	"log"
	"time"
)

// Training - struct for training
type Training struct {
	TableName struct{}  `sql:"trainings"`
	ID        int64     `sql:"id" json:"id" `
	StartDate time.Time `sql:"start_date" json:"start-date"`
	EndDate   time.Time `sql:"end_date" json:"end-date"`
	StartStr  string    `sql:"-" json:"start-str"`
	EndStr    string    `sql:"-" json:"end-str"`
	Notes     string    `sql:"notes, null" json:"notes"`
}

// GetTraining - get training by id
func (e *EDc) GetTraining(id int64) (training Training, err error) {
	if id == 0 {
		return
	}
	err = e.db.Model(&training).Where("id = ?", id).Select()
	if err != nil {
		log.Println("GetTraining ", err)
	}
	return
}

// GetTrainingAll - get all training
func (e *EDc) GetTrainingAll() (trainings []Training, err error) {
	err = e.db.Model(&trainings).Order("name ASC").Select()
	if err != nil {
		log.Println("GetTrainingAll ", err)
		return
	}
	for i := range trainings {
		trainings[i].StartStr = setStrMonth(trainings[i].StartDate)
		trainings[i].EndStr = setStrMonth(trainings[i].EndDate)
	}
	return
}

// CreateTraining - create new training
func (e *EDc) CreateTraining(training Training) (err error) {
	err = e.db.Create(&training)
	if err != nil {
		log.Println("CreateTraining ", err)
	}
	return
}

// UpdateTraining - save changes to training
func (e *EDc) UpdateTraining(training Training) (err error) {
	err = e.db.Update(&training)
	if err != nil {
		log.Println("UpdateTraining ", err)
	}
	return
}

// DeleteTraining - delete training by id
func (e *EDc) DeleteTraining(id int64) (err error) {
	if id == 0 {
		return
	}
	_, err = e.db.Exec("DELETE * FROM trainings WHERE id = ?", id)
	if err != nil {
		log.Println("DeleteTraining ", err)
	}
	return
}

func (e *EDc) trainingCreateTable() (err error) {
	str := `CREATE TABLE IF NOT EXISTS trainings (id bigserial primary key, start_date date, end_date date, notes text)`
	_, err = e.db.Exec(str)
	if err != nil {
		log.Println("trainingCreateTable ", err)
	}
	return
}
