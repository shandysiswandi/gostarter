package grpc

import (
	"context"
	"strconv"

	pb "github.com/shandysiswandi/gostarter/api/gen-proto/todo"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type Endpoint struct {
	pb.UnimplementedTodoServiceServer

	validator validation.Validator

	findUC         domain.Find
	fetchUC        domain.Fetch
	createUC       domain.Create
	deleteUC       domain.Delete
	updateUC       domain.Update
	updateStatusUC domain.UpdateStatus
}

func NewEndpoint(
	validator validation.Validator,
	findUC domain.Find,
	fetchUC domain.Fetch,
	createUC domain.Create,
	deleteUC domain.Delete,
	updateUC domain.Update,
	updateStatusUC domain.UpdateStatus,
) *Endpoint {
	return &Endpoint{
		validator:      validator,
		findUC:         findUC,
		fetchUC:        fetchUC,
		createUC:       createUC,
		deleteUC:       deleteUC,
		updateUC:       updateUC,
		updateStatusUC: updateStatusUC,
	}
}

func (e *Endpoint) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	if err := e.validator.Validate(req); err != nil {
		return nil, err
	}

	resp, err := e.createUC.Execute(ctx, domain.CreateInput{Title: req.GetTitle(), Description: req.GetDescription()})
	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{Id: resp.ID}, nil
}

func (e *Endpoint) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if err := e.validator.Validate(req); err != nil {
		return nil, err
	}

	resp, err := e.deleteUC.Execute(ctx, domain.DeleteInput{ID: req.GetId()})
	if err != nil {
		return nil, err
	}

	return &pb.DeleteResponse{Id: resp.ID}, nil
}

func (e *Endpoint) Find(ctx context.Context, req *pb.FindRequest) (*pb.FindResponse, error) {
	if err := e.validator.Validate(req); err != nil {
		return nil, err
	}

	resp, err := e.findUC.Execute(ctx, domain.FindInput{ID: req.GetId()})
	if err != nil {
		return nil, err
	}

	return &pb.FindResponse{
		Id:          req.GetId(),
		Title:       resp.Title,
		Description: resp.Description,
		Status:      pb.Status(resp.Status),
	}, nil
}

func (e *Endpoint) Fetch(ctx context.Context, req *pb.FetchRequest) (
	*pb.FetchResponse, error,
) {
	if err := e.validator.Validate(req); err != nil {
		return nil, err
	}

	resp, err := e.fetchUC.Execute(ctx, domain.FetchInput{
		ID:          strconv.FormatUint(req.GetId(), 10),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Status:      req.GetStatus().String(),
	})
	if err != nil {
		return nil, err
	}

	todos := make([]*pb.Todo, 0)
	for _, todo := range resp {
		todos = append(todos, &pb.Todo{
			Id:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			Status:      pb.Status(todo.Status),
		})
	}

	return &pb.FetchResponse{Todos: todos}, nil
}

func (e *Endpoint) UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.UpdateStatusResponse, error) {
	if err := e.validator.Validate(req); err != nil {
		return nil, err
	}

	resp, err := e.updateStatusUC.Execute(ctx, domain.UpdateStatusInput{
		ID:     req.GetId(),
		Status: req.GetStatus().String(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.UpdateStatusResponse{Id: resp.ID, Status: pb.Status(resp.Status)}, nil
}

func (e *Endpoint) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	if err := e.validator.Validate(req); err != nil {
		return nil, err
	}

	resp, err := e.updateUC.Execute(ctx, domain.UpdateInput{
		ID:          req.GetId(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Status:      req.GetStatus().String(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.UpdateResponse{
		Id:          resp.ID,
		Title:       resp.Title,
		Description: resp.Description,
		Status:      req.GetStatus(),
	}, nil
}
