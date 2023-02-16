package telegram

import (
	"testing"
	"time"

	district "github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/district/postgres"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage"
)

func Test_buildAnswer(t *testing.T) {
	myTime := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)
	type args struct {
		d district.District
		o []outage.Outage
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"no_outage", 
			args{d:district.District{City:"Irkutsk",Name:"Novolenino"}, o:[] outage.Outage{}}, 
			"\nDid you mean _*Irkutsk Novolenino*_?\n\n*There is no outages planned in your neigborhood in the closest time*\n",false,
		},
		{
			"got_outages", 
			args{d:district.District{City:"Irkutsk",Name:"Novolenino"}, o:[] outage.Outage{
				{Resource:"water", District:"Novolenino",City: "Irkutsk", StartDate:  myTime.Add(3 * time.Hour),Duration:  3*time.Hour, EndDate:  myTime.Add(4 * time.Hour), Notes: "",SourceURL:  "test"},
				{Resource:"water", District:"Novolenino",City: "Irkutsk", StartDate:  myTime.Add(3 * time.Hour),Duration:  3*time.Hour, EndDate:  myTime.Add(4 * time.Hour), Notes: "pushkina, kolotushkina", SourceURL: "test"},
			},			
			}, 
			"\nDid you mean _*Irkutsk Novolenino*_?\n\n*Here are the closest outages found for your neigborhood:*\n\n*water outage* from 2006-01-02 21:04:05 to 2006-01-02 22:04:05\n\n*water outage* from 2006-01-02 21:04:05 to 2006-01-02 22:04:05\n\n*In the next areas and streets:*\npushkina, kolotushkina\n\n",false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildAnswer(tt.args.d, tt.args.o)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildAnswer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("buildAnswer() = %v, want %v", got, tt.want)
			}
		})
	}
}
