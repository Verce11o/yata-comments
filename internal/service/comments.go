package service

import (
	"context"
	"github.com/Verce11o/yata-comments/internal/domain"
	"github.com/Verce11o/yata-comments/internal/lib/grpc_errors"
	"github.com/Verce11o/yata-comments/internal/repository"
	pb "github.com/Verce11o/yata-protos/gen/go/comments"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Comment struct {
	log    *zap.SugaredLogger
	tracer trace.Tracer
	repo   repository.PostgresRepository
	redis  repository.RedisRepository
	minio  repository.MinioRepository
}

func NewCommentService(log *zap.SugaredLogger, tracer trace.Tracer, repo repository.PostgresRepository, redis repository.RedisRepository, minio repository.MinioRepository) *Comment {
	return &Comment{log: log, tracer: tracer, repo: repo, redis: redis, minio: minio}
}

func (t *Comment) CreateComment(ctx context.Context, input *pb.CreateCommentRequest) (string, error) {
	ctx, span := t.tracer.Start(ctx, "commentService.CreateComment")
	defer span.End()

	image := input.GetImage()

	if image != nil {

		err := t.minio.AddCommentImage(ctx, image, image.GetName())

		if err != nil {
			t.log.Errorf("cannot add image to comment in minio: %v", err.Error())
		}

	}

	commentID, err := t.repo.CreateComment(ctx, input, image.GetName())

	if err != nil {
		return "", err
	}

	return commentID, nil
}

func (t *Comment) GetComment(ctx context.Context, commentID string) (domain.Comment, error) {
	ctx, span := t.tracer.Start(ctx, "commentService.GetComment")
	defer span.End()

	cachedComment, err := t.redis.GetCommentByIDCtx(ctx, commentID)

	if err != nil {
		t.log.Infof("cannot get comment by id in redis: %v", err.Error())
	}

	if cachedComment != nil {
		t.log.Info("returned cache")
		return *cachedComment, nil
	}

	comment, err := t.repo.GetComment(ctx, commentID)

	if err != nil {
		t.log.Errorf("cannot get comment by id in postgres: %v", err.Error())
		return domain.Comment{}, err
	}

	if err := t.redis.SetByIDCtx(ctx, commentID, comment); err != nil {
		t.log.Errorf("cannot set comment by id in redis: %v", err.Error())
	}

	return *comment, nil

}

func (t *Comment) GetAllTweetComments(ctx context.Context, input *pb.GetAllTweetCommentsRequest) ([]*pb.Comment, string, error) {
	ctx, span := t.tracer.Start(ctx, "commentService.GetAllComments")
	defer span.End()

	t.log.Debugf("")

	comments, nextCursor, err := t.repo.GetAllTweetComments(ctx, input.GetCursor(), input.GetTweetId())

	if err != nil {
		t.log.Errorf("cannot get all comments by cursor: %v err: %v", input.GetCursor(), err)
	}

	return comments, nextCursor, nil

}

func (t *Comment) UpdateComment(ctx context.Context, input *pb.UpdateCommentRequest) (*domain.Comment, error) {
	ctx, span := t.tracer.Start(ctx, "commentService.UpdateComment")
	defer span.End()

	comment, err := t.repo.GetComment(ctx, input.GetCommentId())

	if err != nil {
		t.log.Errorf("cannot get comment by id in postgres: %v", err.Error())
		return nil, err
	}

	if comment.UserID.String() != input.GetUserId() {
		t.log.Errorf("cannot update comment by id: permission denied")
		return nil, grpc_errors.ErrPermissionDenied
	}

	image := input.GetImage()
	newImageName := comment.ImageName

	if image != nil { // if input image is not nil, we need to update it
		var err error

		if comment.ImageName == "" {
			err = t.minio.AddCommentImage(ctx, image, image.GetName())
		} else {
			err = t.minio.UpdateCommentImage(ctx, comment.ImageName, image.GetName(), image)
		}

		if err != nil {
			t.log.Errorf("cannot update comment image: %v", err.Error())
			return nil, err
		}

		newImageName = image.GetName()
	}

	newComment, err := t.repo.UpdateComment(ctx, input, newImageName)

	if err != nil {
		t.log.Errorf("cannot update comment: %v", err.Error())
		return nil, err
	}

	if err := t.redis.DeleteCommentByIDCtx(ctx, comment.CommentID.String()); err != nil {
		t.log.Errorf("cannot remove comment by id in redis: %v", err.Error())
	}

	return newComment, nil
}

func (t *Comment) DeleteComment(ctx context.Context, input *pb.DeleteCommentRequest) error {
	ctx, span := t.tracer.Start(ctx, "commentService.DeleteComment")
	defer span.End()

	comment, err := t.repo.GetComment(ctx, input.GetCommentId())

	if err != nil {
		t.log.Errorf("cannot get comment by id in postgres: %v", err.Error())
		return err
	}

	if comment.UserID.String() != input.GetUserId() {
		t.log.Errorf("cannot delete comment by id: permission denied")
		return grpc_errors.ErrPermissionDenied
	}

	err = t.repo.DeleteComment(ctx, comment.CommentID.String())

	if err != nil {
		t.log.Errorf("cannot delete comment by id: %v", err.Error())
		return err
	}

	if err := t.redis.DeleteCommentByIDCtx(ctx, comment.CommentID.String()); err != nil {
		t.log.Errorf("cannot delete comment by id in redis: %v", err.Error())
	}

	err = t.minio.DeleteFile(ctx, comment.ImageName)

	if err != nil {
		t.log.Errorf("cannot delete comment image by id: %v", err.Error())
		return err
	}

	return nil

}
