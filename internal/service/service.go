package service

import (
	"btc-project/conf"
	"btc-project/internal/dao"
	"context"
)

// Service service.
type Service struct {
	c   	*conf.Config
	Dao 	*dao.Dao
}

// New new a service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:  c,
		Dao: dao.New(c),
	}

	return s
}

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.Dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
	s.Dao.Close()
}
