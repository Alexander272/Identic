package services

import "context"

type Runner interface {
	Run(ctx context.Context) error
}

// В Services
func (s *Services) GetRunners() []Runner {
	return []Runner{s.OrdersStream} // возвращаем список всех фоновых задач
}
