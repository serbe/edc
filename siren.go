package edc

import (
	"context"
	"time"
)

// Siren - struct for siren
type Siren struct {
	ID          int64  `sql:"id"            json:"id"            form:"id"            query:"id"`
	NumID       int64  `sql:"num_id"        json:"num_id"        form:"num_id"        query:"num_id"`
	NumPass     string `sql:"num_pass"      json:"num_pass"      form:"num_pass"      query:"num_pass"`
	SirenTypeID int64  `sql:"siren_type_id" json:"siren_type_id" form:"siren_type_id" query:"siren_type_id"`
	Address     string `sql:"address"       json:"address"       form:"address"       query:"address"`
	Radio       string `sql:"radio"         json:"radio"         form:"radio"         query:"radio"`
	Desk        string `sql:"desk"          json:"desk"          form:"desk"          query:"desk"`
	ContactID   int64  `sql:"contact_id"    json:"contact_id"    form:"contact_id"    query:"contact_id"`
	CompanyID   int64  `sql:"company_id"    json:"company_id"    form:"company_id"    query:"company_id"`
	Latitude    string `sql:"latitude"      json:"latitude"      form:"latitude"      query:"latitude"`
	Longitude   string `sql:"longitude"     json:"longitude"     form:"longitude"     query:"longitude"`
	Stage       int64  `sql:"stage"         json:"stage"         form:"stage"         query:"stage"`
	Own         string `sql:"own"           json:"own"           form:"own"           query:"own"`
	Note        string `sql:"note"          json:"note"          form:"note"          query:"note"`
	CreatedAt   string `sql:"created_at"    json:"-"`
	UpdatedAt   string `sql:"updated_at"    json:"-"`
}

// SirenList - struct for siren list
type SirenList struct {
	ID            int64    `sql:"id"              json:"id"              form:"id"              query:"id"`
	SirenTypeName string   `sql:"siren_type_name" json:"siren_type_name" form:"siren_type_name" query:"siren_type_name"`
	Address       string   `sql:"address"         json:"address"         form:"address"         query:"address"`
	ContactName   string   `sql:"contact_name"    json:"contact_name"    form:"contact_name"    query:"contact_name"`
	Phones        []string `sql:"phones"          json:"phones"          form:"phones"          query:"phones"          pg:",array"`
}

// SirenGet - get one siren by id
func SirenGet(id int64) (Siren, error) {
	var siren Siren
	if id == 0 {
		return siren, nil
	}
	siren.ID = id
	err := pool.QueryRow(context.Background(), `
		SELECT
			num_id,
			num_pass,
			siren_type_id,
			address,
			radio,
			desk,
			contact_id,
			company_id,
			latitude,
			longitude,
			stage,
			own,
			note,
			created_at,
			updated_at
		FROM
			sirens
		WHERE
			id = $1
	`, id).Scan(&siren.NumID, &siren.NumPass, &siren.SirenTypeID, &siren.Address, &siren.Radio, &siren.Desk, &siren.ContactID, &siren.CompanyID,
		&siren.Latitude, &siren.Longitude, &siren.Stage, &siren.Own, &siren.Note, &siren.CreatedAt, &siren.UpdatedAt)
	if err != nil {
		errmsg("SirenGet QueryRow", err)
	}
	return siren, err
}

// SirenListGet - get all siren for list
func SirenListGet() ([]SirenList, error) {
	var sirens []SirenList
	rows, err := pool.Query(context.Background(), `
		SELECT
			s.id,
			s.address,
			t.name AS siren_type_name,
			c.name AS contact_name,
			array_agg(DISTINCT ph.phone) AS phones
        FROM
			sirens AS s
		LEFT JOIN
			siren_types AS t ON s.type_id = t.id
		LEFT JOIN
			contacts AS c ON s.contact_id = c.id
		LEFT JOIN
			phones AS ph ON s.contact_id = ph.contact_id AND ph.fax = false
		GROUP BY
			s.id,
			t.id,
			c.id
		ORDER BY
			t.name ASC
	`)
	if err != nil {
		errmsg("SirenListGet Query", err)
	}
	for rows.Next() {
		var siren SirenList
		err := rows.Scan(&siren.ID, &siren.Address, &siren.SirenTypeName, &siren.ContactName, &siren.Phones)
		if err != nil {
			errmsg("SirenListGet Scan", err)
			return sirens, err
		}
		sirens = append(sirens, siren)
	}
	return sirens, rows.Err()
}

// SirenInsert - create new siren
func SirenInsert(siren Siren) (int64, error) {
	err := pool.QueryRow(context.Background(), `
		INSERT INTO sirens
		(
			num_id,
			num_pass,
			siren_type_id,
			address,
			radio,
			desk,
			contact_id,
			company_id,
			latitude,
			longitude,
			stage,
			own,
			note,
			created_at,
			updated_at
		)
		VALUES
		(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			$13,
			$14,
			$15
		)
		RETURNING
			id
	`, siren.NumID, siren.NumPass, siren.SirenTypeID, siren.Address, siren.Radio, siren.Desk, siren.ContactID, siren.CompanyID,
		siren.Latitude, siren.Longitude, siren.Stage, siren.Own, siren.Note, time.Now(), time.Now()).Scan(&siren.ID)
	if err != nil {
		errmsg("SirenInsert QueryRow", err)
	}
	return siren.ID, err
}

// SirenUpdate - save siren changes
func SirenUpdate(siren Siren) error {
	_, err := pool.Exec(context.Background(), `
		UPDATE sirens SET
			num_id = $2,
			num_pass = $3,
			siren_type_id = $4,
			address = $5,
			radio = $6,
			desk = $7,
			contact_id = $8,
			company_id = $9,
			latitude = $10,
			longitude = $11,
			stage = $12,
			own = $13,
			note = $14,
			updated_at = $15
		WHERE
			id = $1
	`, siren.ID, siren.NumID, siren.NumPass, siren.SirenTypeID, siren.Address, siren.Radio, siren.Desk, siren.ContactID, siren.CompanyID,
		siren.Latitude, siren.Longitude, siren.Stage, siren.Own, siren.Note, time.Now())
	if err != nil {
		errmsg("SirenUpdate Exec", err)
	}
	return err
}

// SirenDelete - delete siren by id
func SirenDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.Exec(context.Background(), `
		DELETE FROM
			sirens
		WHERE
			id = $1
	`, id)
	if err != nil {
		errmsg("DeleteSiren Exec", err)
	}
	return err
}

func sirenCreateTable() error {
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
	_, err := pool.Exec(context.Background(), str)
	if err != nil {
		errmsg("sirenCreateTable exec", err)
	}
	return err
}
