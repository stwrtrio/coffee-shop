package scheduler

import (
	"context"
	"fmt"
	"time"
)

func (s *scheduler) createMonthlyTables() error {
	now := time.Now()

	for i := 0; i < 2; i++ { // Create tables for the next 2 months
		month := now.AddDate(0, i, 0)
		tableName := fmt.Sprintf("orders_%d%02d", month.Year(), month.Month())

		err := s.OrderRepo.CreateMonthlyOrderTables(context.Background(), tableName)
		if err != nil {
			return err
		}
	}

	return nil
}
