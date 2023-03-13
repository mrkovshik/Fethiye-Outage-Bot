package alert

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage/postgres"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/subscribtion"
	"github.com/pkg/errors"
)

type Alert struct {
	ID             string
	OutageID       string
	SubscribtionID string
	AlertTime      time.Time
	IsSent         bool
	ChatID         int64
}

type AlertRow struct {
	ID             string
	OutageID       string    `db:"outage_id"`
	SubscribtionID string    `db:"subscribtion_id"`
	AlertTime      time.Time `db:"alert_time"`
	IsSent         bool      `db:"is_sent"`
	ChatID         int64
	Cancelled      bool `db:"cancelled"`
}

type AlertStore struct {
	db *sqlx.DB
}

func (from *AlertRow) marshal() Alert {
	return Alert{
		ID:             from.ID,
		OutageID:       from.OutageID,
		SubscribtionID: from.SubscribtionID,
		AlertTime:      from.AlertTime,
		IsSent:         from.IsSent,
		ChatID:         from.ChatID,
	}
}

func (ar *AlertRow) unmarshal(from Alert) error {
	*ar = AlertRow{
		OutageID:       from.OutageID,
		SubscribtionID: from.SubscribtionID,
		AlertTime:      from.AlertTime,
		IsSent:         false,
		Cancelled:      false,
	}
	return nil
}

func NewAlertStore(db *sqlx.DB) *AlertStore {
	return &AlertStore{
		db: db,
	}
}

func (a *AlertStore) GenerateAlertsForNewSub(os postgres.OutageStore, s subscribtion.Subscribtion) error {
	outages, err := os.GetActiveOutagesByCityDistrict(s.DistrictNormalized, s.CityNormalized)
	if err != nil {
		return err
	}
	query := `INSERT INTO alerts (outage_id, subscribtion_id, alert_time, is_sent, cancelled)
	VALUES (:outage_id, :subscribtion_id, :alert_time, :is_sent, :cancelled);`
	for _, o := range outages {
		arow := AlertRow{}
		alert := Alert{
			OutageID:       o.ID,
			SubscribtionID: s.ID,
			AlertTime:      o.StartDate.Add(-time.Duration(s.Period) * time.Hour),
		}
		err := arow.unmarshal(alert)
		if err != nil {
			return err
		}
		if _, err := sqlx.NamedExec(a.db, query, arow); err != nil {
			return errors.Wrap(err, "Error saving alerts")
		}

	}
	return nil
}

func (a *AlertStore) UpdateAlertsOnOutages(os postgres.OutageStore, ss subscribtion.SubscribtionStore) error {
	outages, err := os.GetAllActiveUnalertedOutages()
	if err != nil {
		return err
	}
	query := `INSERT INTO alerts (outage_id, subscribtion_id, alert_time, is_sent, cancelled)
	VALUES (:outage_id, :subscribtion_id, :alert_time, :is_sent, :cancelled);`
	for _, o := range outages {
		subscribtions, err := ss.GetSubsByCityDistrict(o.CityNormalized, o.DistrictNormalized)
		if err != nil {
			return err
		}
		for _, sub := range subscribtions {
			arow := AlertRow{}
			alert := Alert{
				OutageID:       o.ID,
				SubscribtionID: sub.ID,
			}
			err := arow.unmarshal(alert)
			if err != nil {
				return err
			}
			if _, err := sqlx.NamedExec(a.db, query, arow); err != nil {
				return errors.Wrap(err, "Error saving alerts")
			}
		}
		err = os.SetAlerted(o)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AlertStore) CancelByChatID(chatID int64) error {
	query := fmt.Sprintf("UPDATE alerts	SET cancelled = true	WHERE subscribtion_id IN (	SELECT id FROM subscribtions WHERE chat_id = %v) and cancelled=false and is_sent=false;", chatID)
	_, err := a.db.Query(query)
	if err != nil {
		return errors.Wrap(err, "Error deleting alerts")
	}
	return nil
}

func (a *AlertStore) SetIsSent(ID string) error {
	query := fmt.Sprintf("UPDATE alerts  SET is_sent =  true where id='%v'and cancelled=false;", ID)
	_, err := a.db.Query(query)
	if err != nil {
		return errors.Wrap(err, "Error updating alerts")
	}
	return nil
}

func (a *AlertStore) read(query string) ([]Alert, error) {
	rows, err := a.db.Query(query)
	if err != nil {
		err := errors.Wrap(err, "Failed to query database:")
		return []Alert{}, err
	}
	defer rows.Close()
	qryRes := make([]Alert, 0)
	for rows.Next() {
		var a AlertRow
		if err := rows.Scan(&a.ID, &a.ChatID, &a.SubscribtionID, &a.OutageID, &a.AlertTime, &a.IsSent); err != nil {
			err = errors.Wrap(err, "Failed to scan row:")
			return []Alert{}, err
		}
		qryRes = append(qryRes, a.marshal())
	}
	if err := rows.Err(); err != nil {
		err := errors.Wrap(err, "Error iterating through rows:")
		return []Alert{}, err
	}

	return qryRes, err
}

func (a *AlertStore) GetActiveAlerts() ([]Alert, error) {
	query := `SELECT a.id, s.chat_id ,a.subscribtion_id, a.outage_id, a.alert_time, a.is_sent	FROM alerts a	JOIN outages o ON a.outage_id = o.id	JOIN subscribtions s ON a.subscribtion_id  = s.id	WHERE a.alert_time < NOW() AND o.end_date > NOW() and is_sent = false and a.cancelled=false;`
	return a.read(query)
}
