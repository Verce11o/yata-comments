package pagination

import (
	"encoding/base64"
	"fmt"
	"github.com/Verce11o/yata-comments/internal/lib/grpc_errors"
	"github.com/google/uuid"
	"strings"
	"time"
)

func DecodeCursor(encodedCursor string) (time.Time, uuid.UUID, error) {
	byt, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return time.Time{}, [16]byte{}, err
	}

	arrStr := strings.Split(string(byt), ",")
	if len(arrStr) != 2 {
		err = grpc_errors.ErrInvalidCursor
		return time.Time{}, [16]byte{}, err
	}

	res, err := time.Parse(time.RFC3339Nano, arrStr[0])
	if err != nil {
		return time.Time{}, [16]byte{}, err
	}

	tweetID, err := uuid.Parse(arrStr[1])
	if err != nil {
		return time.Time{}, [16]byte{}, err
	}

	return res, tweetID, nil
}

func EncodeCursor(t time.Time, uuid string) string {
	key := fmt.Sprintf("%s,%s", t.Format(time.RFC3339Nano), uuid)
	return base64.StdEncoding.EncodeToString([]byte(key))
}
