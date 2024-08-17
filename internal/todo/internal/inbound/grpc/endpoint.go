package grpc

import (
	"context"
	"strconv"

	pb "github.com/shandysiswandi/gostarter/api/gen-proto/todo"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/usecase"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type Endpoint struct {
	pb.UnimplementedTodoServiceServer

	validator validation.Validator

	getByIDUC       usecase.GetByID
	getWithFilterUC usecase.GetWithFilter
	createUC        usecase.Create
	deleteUC        usecase.Delete
	updateUC        usecase.Update
	updateStatusUC  usecase.UpdateStatus
}

func NewEndpoint(
	validator validation.Validator,
	getByIDUC usecase.GetByID,
	getWithFilterUC usecase.GetWithFilter,
	createUC usecase.Create,
	deleteUC usecase.Delete,
	updateUC usecase.Update,
	updateStatusUC usecase.UpdateStatus,
) *Endpoint {
	return &Endpoint{
		validator:       validator,
		getByIDUC:       getByIDUC,
		getWithFilterUC: getWithFilterUC,
		createUC:        createUC,
		deleteUC:        deleteUC,
		updateUC:        updateUC,
		updateStatusUC:  updateStatusUC,
	}
}

func (e *Endpoint) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	if err := e.validator.Validate(req); err != nil {
		return nil, err
	}

	resp, err := e.createUC.Execute(ctx, usecase.CreateInput{Title: req.GetTitle(), Description: req.GetDescription()})
	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{Id: resp.ID}, nil
}

func (e *Endpoint) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if err := e.validator.Validate(req); err != nil {
		return nil, err
	}

	resp, err := e.deleteUC.Execute(ctx, usecase.DeleteInput{ID: req.GetId()})
	if err != nil {
		return nil, err
	}

	return &pb.DeleteResponse{Id: resp.ID}, nil
}

func (e *Endpoint) GetByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.GetByIDResponse, error) {
	if err := e.validator.Validate(req); err != nil {
		return nil, err
	}

	resp, err := e.getByIDUC.Execute(ctx, usecase.GetByIDInput{ID: req.GetId()})
	if err != nil {
		return nil, err
	}

	return &pb.GetByIDResponse{
		Id:          req.GetId(),
		Title:       resp.Title,
		Description: resp.Description,
		Status:      pb.Status(resp.Status),
	}, nil
}

func (e *Endpoint) GetWithFilter(ctx context.Context, req *pb.GetWithFilterRequest) (
	*pb.GetWithFilterResponse, error,
) {
	if err := e.validator.Validate(req); err != nil {
		return nil, err
	}

	resp, err := e.getWithFilterUC.Execute(ctx, usecase.GetWithFilterInput{
		ID:          strconv.FormatUint(req.GetId(), 10),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Status:      req.GetStatus().String(),
	})
	if err != nil {
		return nil, err
	}

	todos := make([]*pb.Todo, 0)
	for _, todo := range resp.Todos {
		todos = append(todos, &pb.Todo{
			Id:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			Status:      pb.Status(todo.Status),
		})
	}

	return &pb.GetWithFilterResponse{Todos: todos}, nil
}

func (e *Endpoint) UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.UpdateStatusResponse, error) {
	if err := e.validator.Validate(req); err != nil {
		return nil, err
	}

	resp, err := e.updateStatusUC.Execute(ctx, usecase.UpdateStatusInput{
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

	resp, err := e.updateUC.Execute(ctx, usecase.UpdateInput{
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
