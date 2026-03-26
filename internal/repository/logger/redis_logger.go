package logger

import (
	"context"
	"fmt"
	"otus/internal/db"
	"time"
)

const logTTL = 24 * time.Hour

func LogAction(action string, entityType string, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := fmt.Sprintf("log:%s:%s:%d:%d", entityType, action, id, time.Now().UnixNano())
	value := fmt.Sprintf("%s %s id=%d at %s", action, entityType, id, time.Now().Format(time.RFC3339))

	return db.RedisClient.Set(ctx, key, value, logTTL).Err()
}
