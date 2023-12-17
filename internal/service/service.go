package service

import (
	"context"
	"github.com/Verce11o/yata-comments/internal/domain"
	pb "github.com/Verce11o/yata-protos/gen/go/comments"
)

type CommentService interface {
	CreateComment(ctx context.Context, input *pb.CreateCommentRequest) (string, error)
	GetComment(ctx context.Context, commentID string) (domain.Comment, error)
	GetAllTweetComments(ctx context.Context, input *pb.GetAllTweetCommentsRequest) ([]*pb.Comment, string, error)
	UpdateComment(ctx context.Context, input *pb.UpdateCommentRequest) (*domain.Comment, error)
	DeleteComment(ctx context.Context, input *pb.DeleteCommentRequest) error
}
