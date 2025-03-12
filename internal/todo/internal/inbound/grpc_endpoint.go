package inbound

import (
	"context"
	"strings"

	"github.com/shandysiswandi/goreng/telemetry"
	pb "github.com/shandysiswandi/gostarter/api/gen-proto/todo"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
)

type grpcEndpoint struct {
	pb.UnimplementedTodoServiceServer

	tel *telemetry.Telemetry

	findUC         domain.Find
	fetchUC        domain.Fetch
	createUC       domain.Create
	deleteUC       domain.Delete
	updateUC       domain.Update
	updateStatusUC domain.UpdateStatus
}

func (g *grpcEndpoint) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	ctx, span := g.tel.Tracer().Start(ctx, "todo.inbound.grpcEndpoint.Create")
	defer span.End()

	resp, err := g.createUC.Call(ctx, domain.CreateInput{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{Id: resp.ID}, nil
}

func (g *grpcEndpoint) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	ctx, span := g.tel.Tracer().Start(ctx, "todo.inbound.grpcEndpoint.Delete")
	defer span.End()

	resp, err := g.deleteUC.Call(ctx, domain.DeleteInput{ID: req.GetId()})
	if err != nil {
		return nil, err
	}

	return &pb.DeleteResponse{Id: resp.ID}, nil
}

func (g *grpcEndpoint) Find(ctx context.Context, req *pb.FindRequest) (*pb.FindResponse, error) {
	ctx, span := g.tel.Tracer().Start(ctx, "todo.inbound.grpcEndpoint.Find")
	defer span.End()

	resp, err := g.findUC.Call(ctx, domain.FindInput{ID: req.GetId()})
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

func (g *grpcEndpoint) Fetch(ctx context.Context, req *pb.FetchRequest) (*pb.FetchResponse, error) {
	ctx, span := g.tel.Tracer().Start(ctx, "todo.inbound.grpcEndpoint.Fetch")
	defer span.End()

	resp, err := g.fetchUC.Call(ctx, domain.FetchInput{
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

func (g *grpcEndpoint) UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (
	*pb.UpdateStatusResponse, error,
) {
	ctx, span := g.tel.Tracer().Start(ctx, "todo.inbound.grpcEndpoint.UpdateStatus")
	defer span.End()

	resp, err := g.updateStatusUC.Call(ctx, domain.UpdateStatusInput{
		ID:     req.GetId(),
		Status: strings.TrimPrefix(req.GetStatus().String(), "STATUS_"),
	})
	if err != nil {
		return nil, err
	}

	return &pb.UpdateStatusResponse{
		Id:     resp.ID,
		Status: pb.Status(resp.Status.Enum()),
	}, nil
}

func (g *grpcEndpoint) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	ctx, span := g.tel.Tracer().Start(ctx, "todo.inbound.grpcEndpoint.Update")
	defer span.End()

	resp, err := g.updateUC.Call(ctx, domain.UpdateInput{
		ID:          req.GetId(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Status:      strings.TrimPrefix(req.GetStatus().String(), "STATUS_"),
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
