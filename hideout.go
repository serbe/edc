package edc

// Hideout       - защитное сооружение
// ID            - номер в базе данных
// Num           - Номер убежища в реестре имущества
// InvNum        - Инвентарный номер убежища
// InvAdd        - дополнительный код инвентарного номера убежища
// HideoutTypeID - номер типа защитного сооружения в базе данных
// HideoutType   - тип защитного сооружения
// Address       - Полный адрес места расположения убежища, с указанием строения, подъезда
// OwnerID       - номер собственника в базе данных
// Owner         - организация, собственник ЗС
// DesignerID    - номер проектной организации в базе данных
// Designer      - проектная организация
// BuilderID     - номер строительной организации в базе данных
// Builder       - строительная организация
// Purpose       - назначение ЗС в мирное время
// Commissioning - дата ввода в эксплуатацию
// Readiness     - время приведения в готовность
// Capacity      - вместимость
// Area          - общая площадь
// Size          - общий объем
// Floors        - встроено в здание (этажность)
// Separate      - здание отдельно стоящее
// Excavation    - здание в горных выработках
// Inputs        - количество входов
// Coefficient   - коэффициент ослабления гамма излучения К
// Stress        - расчетная нагрузка на действие ударной волны
// Ventilation   - система вентиляции
// Heating       - система отопления
// Power         - система энергосбережения
// Water         - система водоснабжения
// Sewerage      - система канализации
// Implements    - оборудование (инструмент, инвентарь)
// ContactID     - номер контактного лица в базе данных
// Contact       - контактное лицо
// Condition     - готовность к приему укрываемых
// Note          - заметки
// CreatedAt     - время создания записи в базе данных
// UpdatedAt     - время изменения записи в базе данных
type Hideout struct {
	ID            int64       `sql:"id"              json:"id"              form:"id"              query:"id"`
	Num           int64       `sql:"num"             json:"num"             form:"num"             query:"num"`
	InvNum        int64       `sql:"inv_num"         json:"inv_num"         form:"inv_num"         query:"inv_num"`
	InvAdd        int64       `sql:"inv_add"         json:"inv_add"         form:"inv_add"         query:"inv_add"`
	HideoutTypeID int64       `sql:"hideout_type_id" json:"hideout_type_id" form:"hideout_type_id" query:"hideout_type_id"`
	HideoutType   HideoutType `sql:"-"               json:"hideout_type"    form:"hideout_type"    query:"hideout_type"`
	Address       string      `sql:"address"         json:"address"         form:"address"         query:"address"`
	OwnerID       int64       `sql:"owner_id"        json:"owner_id"        form:"owner_id"        query:"owner_id"`
	Owner         Company     `sql:"-"               json:"owner"           form:"owner"           query:"owner"`
	DesignerID    int64       `sql:"designer_id"     json:"designer_id"     form:"designer_id"     query:"designer_id"`
	Designer      Company     `sql:"-"               json:"designer"        form:"designer"        query:"designer"`
	BuilderID     int64       `sql:"builder_id"      json:"builder_id"      form:"builder_id"      query:"builder_id"`
	Builder       Company     `sql:"-"               json:"builder"         form:"builder"         query:"builder"`
	Purpose       string      `sql:"purpose"         json:"purpose"         form:"purpose"         query:"purpose"`
	Commissioning string      `sql:"commissioning"   json:"commissioning"   form:"commissioning"   query:"commissioning"`
	Readiness     int64       `sql:"readiness"       json:"readiness"       form:"readiness"       query:"readiness"`
	Capacity      int64       `sql:"capacity"        json:"capacity"        form:"capacity"        query:"capacity"`
	Area          int64       `sql:"area"            json:"area"            form:"area"            query:"area"`
	Size          int64       `sql:"size"            json:"size"            form:"size"            query:"size"`
	Floors        int64       `sql:"floors"          json:"floors"          form:"floors"          query:"floors"`
	Separate      bool        `sql:"separate"        json:"separate"        form:"separate"        query:"separate"`
	Excavation    bool        `sql:"excavation"      json:"excavation"      form:"excavation"      query:"excavation"`
	Inputs        int64       `sql:"inputs"          json:"inputs"          form:"inputs"          query:"inputs"`
	Coefficient   int64       `sql:"coefficient"     json:"coefficient"     form:"coefficient"     query:"coefficient"`
	Stress        int64       `sql:"stress"          json:"stress"          form:"stress"          query:"stress"`
	Ventilation   string      `sql:"ventilation"     json:"ventilation"     form:"ventilation"     query:"ventilation"`
	Heating       string      `sql:"heating"         json:"heating"         form:"heating"         query:"heating"`
	Power         string      `sql:"power"           json:"power"           form:"power"           query:"power"`
	Water         string      `sql:"water"           json:"water"           form:"water"           query:"water"`
	Sewerage      string      `sql:"sewerage"        json:"sewerage"        form:"sewerage"        query:"sewerage"`
	Implements    string      `sql:"implements"      json:"implements"      form:"implements"      query:"implements"`
	ContactID     int64       `sql:"contact_id"      json:"contact_id"      form:"contact_id"      query:"contact_id"`
	Contact       Contact     `sql:"-"               json:"contact"         form:"contact"         query:"contact"`
	Condition     string      `sql:"condition"       json:"condition"       form:"condition"       query:"condition"`
	Note          string      `sql:"note"            json:"note"            form:"note"            query:"note"`
	CreatedAt     string      `sql:"created_at"      json:"-"`
	UpdatedAt     string      `sql:"updated_at"      json:"-"`
}

// HideoutList - struct for hideout list
type HideoutList struct {
	ID              int64    `sql:"id"                json:"id"                form:"id"                query:"id"`
	HideoutTypeName string   `sql:"hideout_type_name" json:"hideout_type_name" form:"hideout_type_name" query:"hideout_type_name"`
	Address         string   `sql:"address"           json:"address"           form:"address"           query:"address"`
	ContactName     string   `sql:"contact_name"      json:"contact_name"      form:"contact_name"      query:"contact_name"`
	Phones          []string `sql:"phones"            json:"phones"            form:"phones"            query:"phones"            pg:",array"`
}

// GetHideout - get one hideout by id
func GetHideout(id int64) (Hideout, error) {
	var hideout Hideout
	if id == 0 {
		return hideout, nil
	}
	err := pool.Model(&hideout).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetHideout select", err)
	}
	return hideout, err
}

// GetHideoutList - get all hideout for list
func GetHideoutList() ([]HideoutList, error) {
	var hideouts []HideoutList
	_, err := pool.Query(&hideouts, `
		SELECT
			s.id,
			s.address,
			t.name AS hideout_type_name,
			c.name AS contact_name,
			array_agg(DISTINCT ph.phone) AS phones
        FROM
			hideouts AS s
		LEFT JOIN
			hideout_types AS t ON s.type_id = t.id
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
		errmsg("GetHideoutList Query", err)
	}
	return hideouts, err
}

// CreateHideout - create new hideout
func CreateHideout(hideout Hideout) (int64, error) {
	err := pool.Insert(&hideout)
	if err != nil {
		errmsg("CreateHideout insert", err)
	}
	return hideout.ID, err
}

// UpdateHideout - save hideout changes
func UpdateHideout(hideout Hideout) error {
	err := pool.Update(&hideout)
	if err != nil {
		errmsg("UpdateHideout update", err)
	}
	return err
}

// DeleteHideout - delete hideout by id
func DeleteHideout(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.Model(&Hideout{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeleteHideout delete", err)
	}
	return err
}

func hideoutCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			hideouts (
				id              bigserial PRIMARY KEY,
				num             bigint,
				inv_num         bigint,
				inv_add         bigint,
				hideout_type_id bigint,
				address         text,
				owner_id        bigint,
				designer_id     bigint,
				builder_id      bigint,
				purpose         text,
				commissioning   text,
				readiness       bigint,
				capacity        bigint,
				area            bigint,
				size            bigint,
				floors          bigint,
				separate        bool
				excavation      bool
				inputs          bigint,
				coefficient     bigint,
				stress          bigint,
				ventilation     text,
				heating         text,
				power           text,
				water           text,
				sewerage        text,
				implements      text,
				contact_id      bigint,
				condition       text,
				note            text,
				created_at      TIMESTAMP without time zone,
				updated_at      TIMESTAMP without time zone default now(),
				UNIQUE(num, inv_num, inv_add)
			)
	`
	_, err := pool.Exec(str)
	if err != nil {
		errmsg("hideoutCreateTable exec", err)
	}
	return err
}
