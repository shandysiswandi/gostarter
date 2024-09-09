package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/shandysiswandi/gostarter/internal/region/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type SearchStore interface {
	Provinces(ctx context.Context, ids ...string) ([]domain.Province, error)
	Cities(ctx context.Context, parentID string, ids ...string) ([]domain.City, error)
	Districts(ctx context.Context, parentID string, ids ...string) ([]domain.District, error)
	Villages(ctx context.Context, parentID string, ids ...string) ([]domain.Village, error)
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

func (s *Search) Call(ctx context.Context, in domain.SearchInput) ([]domain.Region, error) {
	in.By = strings.TrimSpace(strings.ToLower(in.By))
	in.ParentID = strings.TrimSpace(in.ParentID)
	if err := s.validate.Validate(in); err != nil {
		return nil, goerror.NewInvalidInput("validation input fail", err)
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
		return nil, nil
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

func (s *Search) fromProvinces(ps []domain.Province, err error) ([]domain.Region, error) {
	if err != nil {
		return nil, goerror.NewServer("failed to search provinces", err)
	}

	rs := make([]domain.Region, 0)
	for _, p := range ps {
		rs = append(rs, p.ToRegion())
	}

	return rs, nil
}

func (s *Search) fromCities(cs []domain.City, err error) ([]domain.Region, error) {
	if err != nil {
		return nil, goerror.NewServer("failed to search cities", err)
	}

	rs := make([]domain.Region, 0)
	for _, c := range cs {
		rs = append(rs, c.ToRegion())
	}

	return rs, nil
}

func (s *Search) fromDistricts(ds []domain.District, err error) ([]domain.Region, error) {
	if err != nil {
		return nil, goerror.NewServer("failed to search districts", err)
	}

	rs := make([]domain.Region, 0)
	for _, d := range ds {
		rs = append(rs, d.ToRegion())
	}

	return rs, nil
}

func (s *Search) fromVillages(vs []domain.Village, err error) ([]domain.Region, error) {
	if err != nil {
		return nil, goerror.NewServer("failed to search villages", err)
	}

	rs := make([]domain.Region, 0)
	for _, v := range vs {
		rs = append(rs, v.ToRegion())
	}

	return rs, nil
}
