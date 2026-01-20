package queue

import (
	"shared/database"

	"github.com/hibiken/asynq"
)

func NewAsynqServer(cfg database.AsynqConfig) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     cfg.Addr,
			Password: cfg.Password,
			DB:       cfg.DB,
		},
		asynq.Config{
			Concurrency: 2,
			Queues: map[string]int{
				"reports": 5,
			},
			Logger: &AsynqLogger{},
		},
	)
}
