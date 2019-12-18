package edc

import "context"

// HideoutType - struct for hideoutType
type HideoutType struct {
	ID        int64  `sql:"id"         json:"id"   form:"id"   query:"id"`
	Name      string `sql:"name"       json:"name" form:"name" query:"name"`
	Note      string `sql:"note"       json:"note" form:"note" query:"note"`
	CreatedAt string `sql:"created_at" json:"-"`
	UpdatedAt string `sql:"updated_at" json:"-"`
}

// HideoutTypeList - struct for hideoutType list
type HideoutTypeList struct {
	ID   int64  `sql:"id"   json:"id"   form:"id"   query:"id"`
	Name string `sql:"name" json:"name" form:"name" query:"name"`
	Note string `sql:"note" json:"note" form:"note" query:"note"`
}

// HideoutTypeGet - get one hideoutType by id
func HideoutTypeGet(id int64) (HideoutType, error) {
	var hideoutType HideoutType
	if id == 0 {
		return hideoutType, nil
	}
	err := pool.QueryRow(context.Background(), &hideoutType).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetHideoutType select", err)
	}
	return hideoutType, err
}

// HideoutTypeListGet - get all hideoutType for list
func HideoutTypeListGet() ([]HideoutTypeList, error) {
	var hideoutTypes []HideoutTypeList
	err := pool.QueryRow(context.Background(), &HideoutType{}).
		Column("id", "name", "note").
		Order("name ASC").
		Select(&hideoutTypes)
	if err != nil {
		errmsg("GetHideoutTypeList select", err)
	}
	return hideoutTypes, err
}

// HideoutTypeSelectGet - get all hideoutType for select
func HideoutTypeSelectGet() ([]SelectItem, error) {
	var hideoutTypes []SelectItem
	rows, err := pool.Query(context.Background(), `
		SELECT
			id,
			name
		FROM
			hideout_types
		ORDER BY
			name ASC
	`)
	if err != nil {
		errmsg("HideoutTypeSelectGet Query", err)
	}
	for rows.Next() {
		var hideoutType SelectItem
		err := rows.Scan(&hideoutType.ID, &hideoutType.Name)
		if err != nil {
			errmsg("HideoutTypeSelectGet Scan", err)
			return hideoutTypes, err
		}
		hideoutTypes = append(hideoutTypes, hideoutType)
	}
	return hideoutTypes, rows.Err()
}

// HideoutTypeInsert - create new hideoutType
func HideoutTypeInsert(hideoutType HideoutType) (int64, error) {
	err := pool.Insert(&hideoutType)
	if err != nil {
		errmsg("CreateHideoutType insert", err)
	}
	return hideoutType.ID, nil
}

// HideoutTypeUpdate - save hideoutType changes
func HideoutTypeUpdate(hideoutType HideoutType) error {
	err := pool.Update(&hideoutType)
	if err != nil {
		errmsg("UpdateHideoutType update", err)
	}
	return err
}

// HideoutTypeDelete - delete hideoutType by id
func HideoutTypeDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.QueryRow(context.Background(), &HideoutType{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteHideoutTypedelete", err)
	}
	return err
}

func hideoutTypeCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			hideout_types (
				id         bigserial primary key,
				name       text,
				note       text,
				created_at TIMESTAMP without time zone,
				updated_at
 TIMESTAMP without time zone default now(),
				UNIQUE(name)
			);`
	_, err := pool.Exec(str)
	if err != nil {
		errmsg("hideoutCreateTable exec", err)
	}
	return err
}
