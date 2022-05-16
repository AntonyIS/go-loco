package redis

import (
	"github.com/AntonyIS/go-loco/app"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticache"
)

type redisRepository struct {
	Client    *elasticache.ElastiCache
	Tablename string
}

func newRedisClient() *elasticache.ElastiCache {
	mySession := session.Must(session.NewSession())

	// Create a ElastiCache client from just a session.
	client := elasticache.New(mySession)

	// Create a ElastiCache client with additional configuration
	return client
}

func NewRedisReposistory(tablename string) (app.LocomotiveRepository, error) {
	repo := &redisRepository{
		Client: newRedisClient(),
	}
	return repo, nil
}
