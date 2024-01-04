package postgres

import (
	"context"
	"database/sql"
	"github.com/Verce11o/yata-comments/internal/domain"
	"github.com/Verce11o/yata-comments/internal/lib/pagination"
	pb "github.com/Verce11o/yata-protos/gen/go/comments"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

const (
	paginationLimit = 10
)

type CommentsPostgres struct {
	db     *sqlx.DB
	tracer trace.Tracer
}

func NewCommentsPostgres(db *sqlx.DB, tracer trace.Tracer) *CommentsPostgres {
	return &CommentsPostgres{db: db, tracer: tracer}
}

func (c *CommentsPostgres) CreateComment(ctx context.Context, input *pb.CreateCommentRequest, imageName string) (string, error) {
	ctx, span := c.tracer.Start(ctx, "commentPostgres.CreateTweet")
	defer span.End()

	var commentID string

	q := "INSERT INTO comments (tweet_id, user_id, text, image_name) VALUES ($1, $2, $3, $4) RETURNING comment_id"

	stmt, err := c.db.PreparexContext(ctx, q)

	if err != nil {
		return "", err
	}

	err = stmt.QueryRowxContext(ctx, input.GetTweetId(), input.GetUserId(), input.GetText(), imageName).Scan(&commentID)

	if err != nil {
		return "", err
	}

	return commentID, nil

}

func (c *CommentsPostgres) GetComment(ctx context.Context, CommentID string) (*domain.Comment, error) {
	ctx, span := c.tracer.Start(ctx, "CommentPostgres.GetComment")
	defer span.End()

	var comment domain.Comment

	q := "SELECT * FROM comments WHERE comment_id = $1"

	err := c.db.QueryRowxContext(ctx, q, CommentID).StructScan(&comment)

	if err != nil {
		return nil, sql.ErrNoRows
	}

	return &comment, nil
}

func (c *CommentsPostgres) GetAllTweetComments(ctx context.Context, cursor string, tweetID string) ([]*pb.Comment, string, error) {
	ctx, span := c.tracer.Start(ctx, "commentPostgres.GetAllComments")
	defer span.End()

	var createdAt time.Time
	var commentID uuid.UUID
	var err error

	if cursor != "" {
		createdAt, commentID, err = pagination.DecodeCursor(cursor)
		if err != nil {
			return nil, "", err
		}
	}

	q := "SELECT * FROM comments WHERE (created_at, comment_id) > ($1, $2) AND tweet_id = $3 ORDER BY created_at, comment_id LIMIT $4"

	rows, err := c.db.QueryxContext(ctx, q, createdAt, commentID, tweetID, paginationLimit)

	if err != nil {
		return nil, "", err
	}

	var comments []*pb.Comment
	var latestCreatedAt time.Time

	for rows.Next() {
		var item domain.Comment
		err = rows.StructScan(&item)
		if err != nil {
			return nil, "", err
		}
		comments = append(comments, &pb.Comment{
			TweetId:   tweetID,
			UserId:    item.UserID.String(),
			CommentId: item.CommentID.String(),
			Text:      item.Text,
			CreatedAt: timestamppb.New(item.CreatedAt),
		})
		latestCreatedAt = item.CreatedAt
	}

	var nextCursor string
	if len(comments) > 0 {
		nextCursor = pagination.EncodeCursor(latestCreatedAt, comments[len(comments)-1].CommentId)
	}

	return comments, nextCursor, nil
}

func (c *CommentsPostgres) UpdateComment(ctx context.Context, input *pb.UpdateCommentRequest, imageName string) (*domain.Comment, error) {
	ctx, span := c.tracer.Start(ctx, "commentPostgres.Updatecomment")
	defer span.End()

	var comment domain.Comment

	q := "UPDATE comments SET text = $1, image_name = $2, updated_at = CURRENT_TIMESTAMP WHERE comment_id = $3 RETURNING *"

	if err := c.db.QueryRowxContext(ctx, q, input.GetText(), imageName, input.GetCommentId()).StructScan(&comment); err != nil {
		return nil, err
	}

	return &comment, nil
}

func (c *CommentsPostgres) DeleteComment(ctx context.Context, commentID string) error {
	ctx, span := c.tracer.Start(ctx, "commentPostgres.DeleteComment")
	defer span.End()

	q := "DELETE FROM comments WHERE comment_id = $1"

	res, err := c.db.ExecContext(ctx, q, commentID)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
