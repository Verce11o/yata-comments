package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Verce11o/yata-comments/internal/domain"
	"github.com/Verce11o/yata-comments/internal/lib/grpc_errors"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
	"time"
)

const (
	commentTTL = 3600
)

type CommentsRedis struct {
	client *redis.Client
	tracer trace.Tracer
}

func NewCommentsRedis(client *redis.Client, tracer trace.Tracer) *CommentsRedis {
	return &CommentsRedis{client: client, tracer: tracer}
}

func (r *CommentsRedis) GetCommentByIDCtx(ctx context.Context, commentID string) (*domain.Comment, error) {
	ctx, span := r.tracer.Start(ctx, "commentRedis.GetCommentByIDCtx")
	defer span.End()

	commentBytes, err := r.client.Get(ctx, r.createKey(commentID)).Bytes()

	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, grpc_errors.ErrNotFound
		}
		return nil, err
	}

	var comment domain.Comment

	if err = json.Unmarshal(commentBytes, &comment); err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *CommentsRedis) SetByIDCtx(ctx context.Context, commentID string, comment *domain.Comment) error {
	ctx, span := r.tracer.Start(ctx, "commentRedis.SetByIDCtx")
	defer span.End()

	commentBytes, err := json.Marshal(comment)

	if err != nil {
		return err
	}

	return r.client.Set(ctx, r.createKey(commentID), commentBytes, time.Second*time.Duration(commentTTL)).Err()
}

func (r *CommentsRedis) DeleteCommentByIDCtx(ctx context.Context, commentID string) error {
	ctx, span := r.tracer.Start(ctx, "commentRedis.DeleteCommentByIDCtx")
	defer span.End()

	return r.client.Del(ctx, r.createKey(commentID)).Err()
}

func (r *CommentsRedis) createKey(key string) string {
	return fmt.Sprintf("comment:%s", key)
}
