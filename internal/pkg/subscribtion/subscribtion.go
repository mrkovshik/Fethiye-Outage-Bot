package subscribtion

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/util"

	"github.com/pkg/errors"
)

type Subscribtion struct {
	ID                 string
	ChatID             int64
	City               string
	District           string
	CityNormalized     string
	DistrictNormalized string
	Period             int
}
type subscribtionRow struct {
	ID                 string `db:"id"`
	ChatID             int64  `db:"chat_id"`
	City               string `db:"city"`
	District           string `db:"district"`
	CityNormalized     string `db:"city_normalized"`
	DistrictNormalized string `db:"district_normalized"`
	Period             int    `db:"period"`
	Cancelled          bool   `db:"cancelled"`
}

func (from *subscribtionRow) marshal() Subscribtion {
	return Subscribtion{
		ID:                 from.ID,
		ChatID:             from.ChatID,
		City:               from.City,
		District:           from.District,
		CityNormalized:     from.CityNormalized,
		DistrictNormalized: from.DistrictNormalized,
		Period:             from.Period,
	}
}

func (su *subscribtionRow) unmarshal(from Subscribtion) error {
	c, err := util.Normalize(from.City)
	if err != nil {
		return err
	}

	s, err := util.Normalize(from.District)
	if err != nil {
		return err
	}
	*su = subscribtionRow{
		ChatID:             from.ChatID,
		City:               from.City,
		District:           from.District,
		CityNormalized:     strings.Join(c, " "),
		DistrictNormalized: strings.Join(s, " "),
		Period:             from.Period,
		Cancelled:          false,
	}
	return nil
}

func (su *subscribtionRow) Columns() []string {
	return []string{
		"id", "chat_id", "city", "district", "city_normalized", "district_normalized", "period",
	}
}

type SubscribtionStore struct {
	db *sqlx.DB
}

func NewSubsribtionStore(db *sqlx.DB) *SubscribtionStore {
	return &SubscribtionStore{
		db: db,
	}
}

func (sstr *SubscribtionStore) Save(ss Subscribtion) error {
	// TODO from Columns
	query := `INSERT INTO subscribtions (chat_id, city, district, city_normalized, district_normalized, period, cancelled)
  			VALUES (:chat_id, :city, :district, :city_normalized, :district_normalized, :period, :cancelled);`

	var srow subscribtionRow
	err := srow.unmarshal(ss)
	if err != nil {
		return err
	}

	if _, err := sqlx.NamedExec(sstr.db, query, srow); err != nil {
		return errors.Wrap(err, "Error Saving Subscribtion")
	}

	return nil
}
func (sstr *SubscribtionStore) read(query string) ([]Subscribtion, error) {
	rows, err := sstr.db.Query(query)
	if err != nil {
		err := errors.Wrap(err, "Failed to query database:")
		return []Subscribtion{}, err
	}
	defer rows.Close()

	qryRes := make([]Subscribtion, 0)
	for rows.Next() {
		var s subscribtionRow
		if err := rows.Scan(&s.ID, &s.City, &s.District, &s.CityNormalized, &s.DistrictNormalized, &s.Period); err != nil {
			err = errors.Wrap(err, "Failed to scan row:")
			return []Subscribtion{}, err
		}
		qryRes = append(qryRes, s.marshal())
	}
	if err := rows.Err(); err != nil {
		err := errors.Wrap(err, "Error iterating through rows:")
		return []Subscribtion{}, err
	}

	return qryRes, err
}

func (sstr *SubscribtionStore) modify(query string) error {
	_, err := sstr.db.Query(query)
	if err != nil {
		err := errors.Wrap(err, "Failed to query database:")
		return err
	}
	return nil
}
func (ss *SubscribtionStore) CancelSubscribtion(s Subscribtion) error {
	query := fmt.Sprintf("Update subscribtions set cancelled=true where chat_id=%v and cancelled=false;", s.ChatID)
	return ss.modify(query)
}

func (ss *SubscribtionStore) ModifyPeriod(s Subscribtion) error {
	query := fmt.Sprintf("UPDATE subscribtions  SET period = %v where chat_id=%v and cancelled=false;", s.Period, s.ChatID)
	return ss.modify(query)
}

func (ss *SubscribtionStore) ModifyLocation(s Subscribtion) error {
	query := fmt.Sprintf("UPDATE subscribtions  SET city = '%v', district = '%v', city_normalized = '%v', district_normalized = '%v' where chat_id=%v and cancelled=false;", s.City, s.District, s.CityNormalized, s.DistrictNormalized, s.ChatID)
	return ss.modify(query)
}

func (ss *SubscribtionStore) GetSubsByCityDistrict(city string, distr string) ([]Subscribtion, error) {
	query := fmt.Sprintf("SELECT id, city, district, city_normalized, district_normalized, period FROM subscribtions WHERE city_normalized='%v' AND district_normalized='%v' and cancelled=false;", city, distr)
	return ss.read(query)
}

func (ss *SubscribtionStore) GetSubsByChatID(chatID int64) ([]Subscribtion, error) {
	query := fmt.Sprintf("SELECT id, city, district, city_normalized, district_normalized, period FROM subscribtions WHERE chat_id=%v and cancelled=false;", chatID)
	return ss.read(query)
}

func (ss *SubscribtionStore) SubExists(id int64) (bool, error) {
	subs, err := ss.GetSubsByChatID(id)
	if err != nil {
		err := errors.Wrap(err, "Failed to query database:")
		return false, err
	}
	if len(subs) == 0 {
		return false, err
	} else {
		return true, err
	}
}
