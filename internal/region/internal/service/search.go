package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/shandysiswandi/gostarter/internal/region/internal/entity"
	"github.com/shandysiswandi/gostarter/internal/region/internal/usecase"
	gerr "github.com/shandysiswandi/gostarter/pkg/gerr"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type SearchStore interface {
	Provinces(ctx context.Context, ids ...string) ([]entity.Province, error)
	Cities(ctx context.Context, parentID string, ids ...string) ([]entity.City, error)
	Districts(ctx context.Context, parentID string, ids ...string) ([]entity.District, error)
	Villages(ctx context.Context, parentID string, ids ...string) ([]entity.Village, error)
}

type Search struct {
	validate validation.Validator
	store    SearchStore
}

func NewSearch(validate validation.Validator, store SearchStore) *Search {
	return &Search{
		validate: validate,
		store:    store,
	}
}

func (s *Search) Execute(ctx context.Context, in usecase.SearchInput) (*usecase.SearchOutput, error) {
	in.By = strings.TrimSpace(strings.ToLower(in.By))
	in.ParentID = strings.TrimSpace(in.ParentID)
	if err := s.validate.Validate(in); err != nil {
		return nil, gerr.NewValidation("validation input fail", err)
	}

	ids := s.parseIDs(in.IDs)

	switch in.By {
	case "provinces":
		resp, err := s.store.Provinces(ctx, ids...)

		return s.fromProvinces(resp, err)
	case "cities":
		resp, err := s.store.Cities(ctx, in.ParentID, ids...)

		return s.fromCities(resp, err)
	case "districts":
		resp, err := s.store.Districts(ctx, in.ParentID, ids...)

		return s.fromDistricts(resp, err)
	case "villages":
		resp, err := s.store.Villages(ctx, in.ParentID, ids...)

		return s.fromVillages(resp, err)
	default:
		return &usecase.SearchOutput{}, nil
	}
}

func (s *Search) parseIDs(ids string) []string {
	result := make([]string, 0)
	idList := strings.Split(ids, ",")
	for _, id := range idList {
		if _, err := strconv.Atoi(id); err == nil {
			result = append(result, id)
		}
	}

	return result
}

func (s *Search) fromProvinces(ps []entity.Province, err error) (*usecase.SearchOutput, error) {
	if err != nil {
		return nil, gerr.NewServer("failed to search provinces", err)
	}

	rs := make([]entity.Region, 0)
	for _, p := range ps {
		rs = append(rs, p.ToRegion())
	}

	return &usecase.SearchOutput{Regions: rs}, nil
}

func (s *Search) fromCities(cs []entity.City, err error) (*usecase.SearchOutput, error) {
	if err != nil {
		return nil, gerr.NewServer("failed to search cities", err)
	}

	rs := make([]entity.Region, 0)
	for _, c := range cs {
		rs = append(rs, c.ToRegion())
	}

	return &usecase.SearchOutput{Regions: rs}, nil
}

func (s *Search) fromDistricts(ds []entity.District, err error) (*usecase.SearchOutput, error) {
	if err != nil {
		return nil, gerr.NewServer("failed to search districts", err)
	}

	rs := make([]entity.Region, 0)
	for _, d := range ds {
		rs = append(rs, d.ToRegion())
	}

	return &usecase.SearchOutput{Regions: rs}, nil
}

func (s *Search) fromVillages(vs []entity.Village, err error) (*usecase.SearchOutput, error) {
	if err != nil {
		return nil, gerr.NewServer("failed to search villages", err)
	}

	rs := make([]entity.Region, 0)
	for _, v := range vs {
		rs = append(rs, v.ToRegion())
	}

	return &usecase.SearchOutput{Regions: rs}, nil
}
