package grpc

import (
	"context"
	"github.com/Verce11o/yata-comments/internal/lib/grpc_errors"
	"github.com/Verce11o/yata-comments/internal/service"
	pb "github.com/Verce11o/yata-protos/gen/go/comments"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

type CommentGRPC struct {
	log     *zap.SugaredLogger
	tracer  trace.Tracer
	service service.CommentService
	pb.UnimplementedCommentsServer
}

func NewCommentGRPC(log *zap.SugaredLogger, tracer trace.Tracer, service service.CommentService) *CommentGRPC {
	return &CommentGRPC{log: log, tracer: tracer, service: service}
}

func (c *CommentGRPC) CreateComment(ctx context.Context, input *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error) {
	ctx, span := c.tracer.Start(ctx, "CreateComment")
	defer span.End()

	commentID, err := c.service.CreateComment(ctx, input)

	if err != nil {
		c.log.Errorf("CreateComment: %v", err.Error())
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "CreateComment: %v", err)
	}

	return &pb.CreateCommentResponse{CommentId: commentID}, nil

}

func (c *CommentGRPC) GetComment(ctx context.Context, input *pb.GetCommentRequest) (*pb.Comment, error) {
	ctx, span := c.tracer.Start(ctx, "GetComment")
	defer span.End()

	comment, err := c.service.GetComment(ctx, input.GetCommentId())

	if err != nil {
		c.log.Errorf("GetComment: %v", err.Error())
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "GetComment: %v", err)
	}

	return &pb.Comment{
		TweetId:   comment.TweetID.String(),
		UserId:    comment.UserID.String(),
		CommentId: comment.CommentID.String(),
		Text:      comment.Text,
		ImageUrl:  &comment.ImageURL,
	}, nil
}

func (c *CommentGRPC) GetAllTweetComments(ctx context.Context, input *pb.GetAllTweetCommentsRequest) (*pb.GetAllCommentsResponse, error) {
	ctx, span := c.tracer.Start(ctx, "GetAllTweetComments")
	defer span.End()

	comments, nextCursor, err := c.service.GetAllTweetComments(ctx, input)

	if err != nil {
		c.log.Errorf("GetAllComments: %v", err.Error())
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "GetAllComments: %v", err)
	}

	return &pb.GetAllCommentsResponse{Comments: comments, Cursor: nextCursor}, nil
}

func (c *CommentGRPC) UpdateComment(ctx context.Context, input *pb.UpdateCommentRequest) (*pb.Comment, error) {
	ctx, span := c.tracer.Start(ctx, "UpdateComment")
	defer span.End()

	comment, err := c.service.UpdateComment(ctx, input)

	if err != nil {
		c.log.Errorf("UpdateComment: %v", err.Error())
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "UpdateComment: %v", err)
	}

	return &pb.Comment{
		TweetId:   comment.TweetID.String(),
		UserId:    comment.UserID.String(),
		CommentId: comment.CommentID.String(),
		Text:      comment.Text,
		ImageUrl:  &comment.ImageURL,
	}, nil
}

func (c *CommentGRPC) DeleteComment(ctx context.Context, input *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	ctx, span := c.tracer.Start(ctx, "DeleteComment")
	defer span.End()

	err := c.service.DeleteComment(ctx, input)

	if err != nil {
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "DeleteComment: %v", err)
	}

	return &pb.DeleteCommentResponse{}, nil
}
