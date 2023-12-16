package repository

import (
	"context"
	"github.com/Verce11o/yata-comments/internal/domain"
	pb "github.com/Verce11o/yata-protos/gen/go/comments"
)

type RedisRepository interface { // maybe refactor ?
	GetCommentByIDCtx(ctx context.Context, key string) (*domain.Comment, error)
	SetByIDCtx(ctx context.Context, commentID string, comment *domain.Comment) error
	DeleteCommentByIDCtx(ctx context.Context, commentID string) error
}

type PostgresRepository interface {
	CreateComment(ctx context.Context, input *pb.CreateCommentRequest, imageName string, imageURL string) (string, error)
	GetComment(ctx context.Context, CommentID string) (*domain.Comment, error)
	GetAllTweetComments(ctx context.Context, cursor string, tweetID string) ([]*pb.Comment, string, error)
	UpdateComment(ctx context.Context, input *pb.UpdateCommentRequest, imageName string, imageURL string) (*domain.Comment, error)
	DeleteComment(ctx context.Context, CommentID string) error
}

type MinioRepository interface {
	AddCommentImage(ctx context.Context, image *pb.Image, fileName string) (string, error)
	GetCommentImage(ctx context.Context, fileName string) (string, error)
	UpdateCommentImage(ctx context.Context, oldName string, newName string, image *pb.Image) error
	DeleteFile(ctx context.Context, fileName string) error
}
