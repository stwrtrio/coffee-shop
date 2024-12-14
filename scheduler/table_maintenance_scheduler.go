package scheduler

import (
	"context"
	"time"
)

func (s *scheduler) createMonthlyTables() error {
	now := time.Now()

	for i := 0; i < 2; i++ { // Create tables for the next 2 months
		timeAdd := now.AddDate(0, i, 0)
		err := s.OrderRepo.CreateMonthlyOrderTables(context.Background(), timeAdd)
		if err != nil {
			return err
		}
	}

	return nil
}
