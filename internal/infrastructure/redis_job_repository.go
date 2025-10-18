package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"go-jobqueue/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisJobRepository struct {
	client *redis.Client
}

func NewRedisJobRepository(client *redis.Client) *RedisJobRepository {
	return &RedisJobRepository{client: client}
}

func (r *RedisJobRepository) Enqueue(job *domain.Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return r.client.LPush(context.Background(), "jobqueue:jobs", data).Err()
}

func (r *RedisJobRepository) Dequeue() (*domain.Job, error) {
	ctx := context.Background()
	data, err := r.client.RPop(ctx, "jobqueue:jobs").Result()
	if err == redis.Nil {
		return nil, nil // queue empty
	}
	if err != nil {
		return nil, err
	}

	var job domain.Job
	if err := json.Unmarshal([]byte(data), &job); err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *RedisJobRepository) UpdateStatus(id string, status domain.JobStatus) error {
	ctx := context.Background()
	val, err := r.client.Get(ctx, "jobqueue:job:"+id).Result()
	if err == redis.Nil {
		return fmt.Errorf("job not found")
	}
	if err != nil {
		return err
	}

	var job domain.Job
	if err := json.Unmarshal([]byte(val), &job); err != nil {
		return err
	}

	job.Status = status
	job.UpdatedAt = time.Now()

	newVal, _ := json.Marshal(job)
	return r.client.Set(ctx, "jobqueue:job:"+id, newVal, 0).Err()
}
