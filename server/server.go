package server

import (
	"context"

	"github.com/entigolabs/waypoint/internal/db"
)

//go:generate go tool oapi-codegen --config=../openapi/oapi-config.yaml -o oapigen.go ../openapi/openapi.yaml

type Server struct {
	db dbReader
}

var _ StrictServerInterface = (*Server)(nil)

func NewServer(database *db.DB) Server {
	return Server{db: database}
}

func (s Server) GetCoreCategories(ctx context.Context, _ GetCoreCategoriesRequestObject) (GetCoreCategoriesResponseObject, error) {
	data, err := listCategories(ctx, s.db)
	if err != nil {
		return nil, err
	}
	return GetCoreCategories200JSONResponse(CategoriesResponse{
		Data:     data,
		Metadata: Metadata{Total: len(data)},
	}), nil
}

func (s Server) GetCoreEmsCategories(ctx context.Context, _ GetCoreEmsCategoriesRequestObject) (GetCoreEmsCategoriesResponseObject, error) {
	data, err := listEmsCategories(ctx, s.db)
	if err != nil {
		return nil, err
	}
	return GetCoreEmsCategories200JSONResponse(EmsCategoriesResponse{
		Data:     data,
		Metadata: Metadata{Total: len(data)},
	}), nil
}

func (s Server) GetCoreEmsThemes(ctx context.Context, _ GetCoreEmsThemesRequestObject) (GetCoreEmsThemesResponseObject, error) {
	data, err := listEmsThemes(ctx, s.db)
	if err != nil {
		return nil, err
	}
	return GetCoreEmsThemes200JSONResponse(EmsThemesResponse{
		Data:     data,
		Metadata: Metadata{Total: len(data)},
	}), nil
}
