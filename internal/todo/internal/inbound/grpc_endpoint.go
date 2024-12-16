package inbound

import (
	"context"

	pb "github.com/shandysiswandi/gostarter/api/gen-proto/todo"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
)

type grpcEndpoint struct {
	pb.UnimplementedTodoServiceServer

	findUC         domain.Find
	fetchUC        domain.Fetch
	createUC       domain.Create
	deleteUC       domain.Delete
	updateUC       domain.Update
	updateStatusUC domain.UpdateStatus
}

func (e *grpcEndpoint) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	resp, err := e.createUC.Call(ctx, domain.CreateInput{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{Id: resp.ID}, nil
}

func (e *grpcEndpoint) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	resp, err := e.deleteUC.Call(ctx, domain.DeleteInput{ID: req.GetId()})
	if err != nil {
		return nil, err
	}

	return &pb.DeleteResponse{Id: resp.ID}, nil
}

func (e *grpcEndpoint) Find(ctx context.Context, req *pb.FindRequest) (*pb.FindResponse, error) {
	resp, err := e.findUC.Call(ctx, domain.FindInput{ID: req.GetId()})
	if err != nil {
		return nil, err
	}

	return &pb.FindResponse{
		Id:          req.GetId(),
		UserId:      resp.UserID,
		Title:       resp.Title,
		Description: resp.Description,
		Status:      pb.Status(resp.Status.Enum()),
	}, nil
}

func (e *grpcEndpoint) Fetch(ctx context.Context, req *pb.FetchRequest) (*pb.FetchResponse, error) {
	resp, err := e.fetchUC.Call(ctx, domain.FetchInput{
		Cursor: req.GetCursor(),
		Limit:  req.GetLimit(),
		Status: req.GetStatus().String(),
	})
	if err != nil {
		return nil, err
	}

	todos := make([]*pb.Todo, 0)
	for _, todo := range resp.Todos {
		todos = append(todos, &pb.Todo{
			Id:          todo.ID,
			UserId:      todo.UserID,
			Title:       todo.Title,
			Description: todo.Description,
			Status:      pb.Status(todo.Status.Enum()),
		})
	}

	return &pb.FetchResponse{
		Todos: todos,
		Pagination: &pb.Pagination{
			NextCursor: resp.NextCursor,
			HasMore:    resp.HasMore,
		},
	}, nil
}

func (e *grpcEndpoint) UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (
	*pb.UpdateStatusResponse, error,
) {
	resp, err := e.updateStatusUC.Call(ctx, domain.UpdateStatusInput{
		ID:     req.GetId(),
		Status: req.GetStatus().String(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.UpdateStatusResponse{
		Id:     resp.ID,
		Status: pb.Status(resp.Status.Enum()),
	}, nil
}

func (e *grpcEndpoint) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	resp, err := e.updateUC.Call(ctx, domain.UpdateInput{
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
		UserId:      resp.UserID,
		Title:       resp.Title,
		Description: resp.Description,
		Status:      pb.Status(resp.Status.Enum()),
	}, nil
}
