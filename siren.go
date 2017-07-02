package edc

// Siren - struct for siren
type Siren struct {
	ID          int64     `sql:"id"              json:"id"            form:"id"            query:"id"`
	NumID       int64     `sql:"num_id,null"     json:"num_id"        form:"num_id"        query:"num_id"`
	NumPass     string    `sql:"num_pass,null"   json:"num_pass"      form:"num_pass"      query:"num_pass"`
	SyrenTypeID int64     `sql:"type_id"         json:"siren_type_id" form:"siren_type_id" query:"siren_type_id"`
	SyrenType   SirenType `sql:"-"               json:"siren_type"    form:"siren_type"    query:"siren_type"`
	Address     string    `sql:"address,null"    json:"address"       form:"address"       query:"address"`
	Radio       string    `sql:"radio,null"      json:"radio"         form:"radio"         query:"radio"`
	Desk        string    `sql:"desk,null"       json:"desk"          form:"desk"          query:"desk"`
	ContactID   int64     `sql:"contact_id,null" json:"contact_id"    form:"contact_id"    query:"contact_id"`
	Contact     Contact   `sql:"-"               json:"contact"       form:"contact"       query:"contact"`
	CompanyID   int64     `sql:"company_id,null" json:"company_id"    form:"company_id"    query:"company_id"`
	Company     Company   `sql:"-"               json:"company"       form:"company"       query:"company"`
	Latitude    string    `sql:"latitude,null"   json:"latitude"      form:"latitude"      query:"latitude"`
	Longitude   string    `sql:"longitude,null"  json:"longitude"     form:"longitude"     query:"longitude"`
	Stage       int64     `sql:"stage,null"      json:"stage"         form:"stage"         query:"stage"`
	Own         string    `sql:"own,null"        json:"own"           form:"own"           query:"own"`
	Note        string    `sql:"note,null"       json:"note"          form:"note"          query:"note"`
	CreatedAt   string    `sql:"created_at"      json:"-"`
	UpdatedAt   string    `sql:"updated_at"      json:"-"`
}

// SirenList - struct for siren list
type SirenList struct {
	ID            int64  `sql:"id"              json:"id"              form:"id"              query:"id"`
	SyrenTypeID   int64  `sql:"type_id"         json:"siren_type_id"   form:"siren_type_id"   query:"siren_type_id"`
	SyrenTypeName int64  `sql:"type_name"       json:"siren_type_name" form:"siren_type_name" query:"siren_type_name"`
	Address       string `sql:"address,null"    json:"address"         form:"address"         query:"address"`
	Note          string `sql:"note,null"       json:"note"            form:"note"            query:"note"`
}

// GetSiren - get one siren by id
func (e *Edb) GetSiren(id int64) (Siren, error) {
	var siren Siren
	if id == 0 {
		return siren, nil
	}
	err := e.db.Model(&siren).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetSiren select", err)
	}
	return siren, err
}

// GetSirenList - get all siren for list
func (e *Edb) GetSirenList() ([]SirenList, error) {
	var sirens []SirenList
	_, err := e.db.Query(&sirens, `
		SELECT
			s.id,
			s.type_id,
			s.address,
			t.name AS type_name,
			s.note
        FROM
			sirens AS s
		LEFT JOIN
			siren_types AS t ON s.type_id = t.id
		ORDER BY
			t.name ASC
	`)
	if err != nil {
		errmsg("GetSirenList Query", err)
	}
	return sirens, err
}

// CreateSiren - create new siren
func (e *Edb) CreateSiren(siren Siren) (int64, error) {
	err := e.db.Insert(&siren)
	if err != nil {
		errmsg("CreateSiren insert", err)
	}
	return siren.ID, err
}

// UpdateSiren - save siren changes
func (e *Edb) UpdateSiren(siren Siren) error {
	err := e.db.Update(&siren)
	if err != nil {
		errmsg("UpdateSiren update", err)
	}
	return err
}

// DeleteSiren - delete siren by id
func (e *Edb) DeleteSiren(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Model(&Siren{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteSiren delete", err)
	}
	return err
}

func (e *Edb) sirenCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			sirens (
				id         bigserial PRIMARY KEY,
				num_id     bigint,
				num_pass   text,
				type_id    bigint,
				address    text,
				radio      text,
				desk       text,
				contact_id bigint,
				company_id bigint,
				latitude   text,
				longitude  text,
				stage      bigint,
				own        text,
				note        text,
				created_at TIMESTAMP without time zone,
				updated_at TIMESTAMP without time zone default now(),
				UNIQUE(num_id, num_pass, type_id)
			)
	`
	_, err := e.db.Exec(str)
	if err != nil {
		errmsg("sirenCreateTable exec", err)
	}
	return err
}
