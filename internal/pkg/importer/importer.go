package importer

import (
	"gitlab.com/v.rianov/favs-backend/internal/pkg/googlesheets"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/maps"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/places"
)

type Importer struct {
	linkResolver maps.LocationLinkResolver
	parser       googlesheets.SheetParser
	service      places.Usecase
}

func NewImporter(linkResolver maps.LocationLinkResolver, parser googlesheets.SheetParser,
	service places.Usecase) Importer {
	return Importer{
		linkResolver: linkResolver,
		parser:       parser,
		service:      service,
	}
}
