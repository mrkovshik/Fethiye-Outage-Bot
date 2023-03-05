package subscribtion

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/pkg/errors"
)

type Subscribtion struct {
	ChatID   int64
	City     string
	District string
	Period   int
}
type subscribtionRow struct {
	ChatID   int64 `db:"chat_id"`
	City     string `db:"city"`
	District string `db:"district"`
	Period   int    `db:"period"`
}

func (su *subscribtionRow) marshal() Subscribtion {
	return Subscribtion{
		ChatID:   su.ChatID,
		City:     su.City,
		District: su.District,
		Period:   su.Period,
	}
}

func (su *subscribtionRow) unmarshal(from Subscribtion) error {
	*su = subscribtionRow{
		ChatID:   from.ChatID,
		City:     from.City,
		District: from.District,
		Period:   from.Period,
	}
	return nil
}

func (su *subscribtionRow) Columns() []string {
	return []string{
		"id", "chat_id", "city", "district", "period",
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
	query := `INSERT INTO subscribtions (chat_id, city, district, period)
  			VALUES (:chat_id, :city, :district, :period);`

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
		if err := rows.Scan(&s.City, &s.District, &s.Period); err != nil {
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
func (ss *SubscribtionStore) RemoveSubscribtion (s Subscribtion ) error {
	query := fmt.Sprintf("Delete from subscribtions where chat_id=%v;", s.ChatID)
	return ss.modify(query)
}

func (ss *SubscribtionStore) ModifyPeriod(s Subscribtion) error {
	query := fmt.Sprintf("UPDATE subscribtions  SET period = %v where chat_id=%v;", s.Period , s.ChatID)
	return ss.modify(query)
}

func (ss *SubscribtionStore) ModifyLocation(s Subscribtion) error {
	query := fmt.Sprintf("UPDATE subscribtions  SET city = '%v', district = '%v' where chat_id=%v;",s.City, s.District, s.ChatID)
	return ss.modify(query)
}

func (ss *SubscribtionStore) GetSubsByDistrict(distr string) ([]Subscribtion, error) {
	query := fmt.Sprintf("SELECT city, district, period FROM subscribtions WHERE district='%v';", distr)
	return ss.read(query)
}

func (ss *SubscribtionStore) GetSubsByChatID(chatID int64) ([]Subscribtion, error) {
	query := fmt.Sprintf("SELECT city, district, period FROM subscribtions WHERE chat_id=%v;", chatID)
	return ss.read(query)
}

func (ss *SubscribtionStore) SubExists (id int64) (bool,error){
	subs,err:=ss.GetSubsByChatID(id)
	if err != nil {
		err := errors.Wrap(err, "Failed to query database:")
		return false, err
	}
	if len(subs)==0{
		return false,err
	} else {
return true,err
	}
}

