package subsribtion

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/subsribtion"
	"github.com/pkg/errors"
)

type subscribtionRow struct {
	ChatID   int    `db:"chat_id"`
	City     string `db:"city"`
	District string `db:"district"`
	UserName string `db:"user_name"`
}

func (su *subscribtionRow) marshal() subsribtion.Subscribtion {
	return subsribtion.Subscribtion{
		ChatID:   su.ChatID,
		City:     su.City,
		District: su.District,
		UserName: su.UserName,
	}
}

func (su *subscribtionRow) unmarshal(from subsribtion.Subscribtion) error {
	*su = subscribtionRow{
		ChatID:   from.ChatID,
		City:     from.City,
		District: from.District,
		UserName: from.UserName,
	}
	return nil
}

func (su *subscribtionRow) Columns() []string {
	return []string{
		"id", "chat_id", "city", "district", "user_name",
	}
}

type SubsribtionStore struct {
	db *sqlx.DB
}

func NewSubsribtionStore(db *sqlx.DB) *SubsribtionStore {
	return &SubsribtionStore{
		db: db,
	}
}

func (sstr *SubsribtionStore) Save(ss subsribtion.Subscribtion) error {
	// TODO from Columns
	query := `INSERT INTO subscribtions (chat_id, city, district, user_name)
  			VALUES (:chat_id, :city, :district, :user_name);`

	var srow subscribtionRow
	err := srow.unmarshal(ss)
	if err != nil {
		return err
	}

	if _, err := sqlx.NamedExec(sstr.db, query, srow); err != nil {
		return errors.Wrap(err, "could not save New synthetic order")
	}

	return nil
}
func (sstr *SubsribtionStore) Read(query string) ([]subsribtion.Subscribtion, error) {
	rows, err := sstr.db.Query(query)
	if err != nil {
		fmt.Println("Failed to query database:", err)
		return []subsribtion.Subscribtion{}, err
	}
	defer rows.Close()

	qryRes := make([]subsribtion.Subscribtion, 0)
	for rows.Next() {
		var s subscribtionRow
		if err := rows.Scan(&s.ChatID, &s.City, &s.District, &s.UserName); err != nil {
			fmt.Println("Failed to scan row:", err)
			return []subsribtion.Subscribtion{}, err
		}
		qryRes = append(qryRes, s.marshal())
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating through rows:", err)
		return []subsribtion.Subscribtion{}, err
	}

	return qryRes, err
}

func (ss *SubsribtionStore) GetSubsByDistrict(distr string) ([]subsribtion.Subscribtion, error) {
	query := fmt.Sprintf("SELECT (chat_id, city, district, user_name) FROM subscribtions WHERE district=%v;", distr)
	return ss.Read(query)
}

func (ss *SubsribtionStore) GetSubsByChatID(chatID int) ([]subsribtion.Subscribtion, error) {
	query := fmt.Sprintf("SELECT (chat_id, city, district, user_name) FROM subscribtions WHERE chat_id=%v;", chatID)
	return ss.Read(query)
}
