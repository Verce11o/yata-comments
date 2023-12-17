package minio

import (
	"bytes"
	"context"
	pb "github.com/Verce11o/yata-protos/gen/go/comments"
	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/otel/trace"
	"time"
)

const (
	userCommentsName = "user-comments"
	imageExpireTime  = time.Hour * 24
)

type CommentMinio struct {
	minio  *minio.Client
	tracer trace.Tracer
}

func NewCommentMinio(minio *minio.Client, tracer trace.Tracer) *CommentMinio {
	return &CommentMinio{minio: minio, tracer: tracer}
}

func (t *CommentMinio) AddCommentImage(ctx context.Context, image *pb.Image, fileName string) error {
	ctx, span := t.tracer.Start(ctx, "commentMinio.AddImage")
	defer span.End()

	reader := bytes.NewReader(image.GetChunk())

	_, err := t.minio.PutObject(
		ctx,
		userCommentsName,
		fileName,
		reader,
		reader.Size(),
		minio.PutObjectOptions{ContentType: image.GetContentType()},
	)
	if err != nil {
		return err
	}

	return nil
}

func (t *CommentMinio) UpdateCommentImage(ctx context.Context, oldName string, newName string, image *pb.Image) error {
	ctx, span := t.tracer.Start(ctx, "commentMinio.UpdateCommentImage")
	defer span.End()

	err := t.DeleteFile(ctx, oldName)

	if err != nil {
		return err
	}

	err = t.AddCommentImage(ctx, image, newName)

	if err != nil {
		return err
	}

	return nil
}

func (t *CommentMinio) DeleteFile(ctx context.Context, fileName string) error {
	ctx, span := t.tracer.Start(ctx, "commentMinio.DeleteFile")
	defer span.End()

	if err := t.minio.RemoveObject(ctx, userCommentsName, fileName, minio.RemoveObjectOptions{}); err != nil {
		return err
	}

	return nil
}
