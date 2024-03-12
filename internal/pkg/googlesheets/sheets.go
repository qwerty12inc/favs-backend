package googlesheets

import (
	"context"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"google.golang.org/api/sheets/v4"
	"log"
	"strings"
)

type SheetParser interface {
	ParsePlaces(ctx context.Context, readRange string) ([]models.GoogleSheetPlace, models.Status)
}

type SheetsParser struct {
	Cl       *sheets.Service
	SheetRef string
}

func NewSheetsParser(cl *sheets.Service, sheetRef string) SheetsParser {
	return SheetsParser{
		Cl:       cl,
		SheetRef: sheetRef,
	}
}

func (s SheetsParser) ParsePlaces(ctx context.Context, readRange string) ([]models.GoogleSheetPlace, models.Status) {
	resp, err := s.Cl.Spreadsheets.Values.Get(s.SheetRef, readRange).Do()
	if err != nil {
		log.Println("Error while reading from sheet: ", err, " range: ", readRange, " sheet: ", s.SheetRef)
		return nil, models.Status{
			Code:    models.InternalError,
			Message: err.Error(),
		}
	}

	values := resp.Values
	places := make([]models.GoogleSheetPlace, 0, len(values))
	for _, row := range values {
		place := models.GoogleSheetPlace{}
		if len(row) > 0 {
			place.Name = row[0].(string)
		}
		if len(row) > 1 {
			labels := strings.Split(row[1].(string), "/")
			for i, label := range labels {
				labels[i] = strings.TrimSpace(label)
			}
			place.Labels = labels
		}
		if len(row) > 2 {
			place.LocationURL = row[2].(string)
		}
		if len(row) > 3 {
			place.Description = row[3].(string)
		}
		if len(row) > 4 {
			place.Instagram = row[4].(string)
		}
		if len(row) > 5 {
			place.Website = row[5].(string)
		}
		places = append(places, place)
	}
	return places, models.Status{Code: models.OK}
}
