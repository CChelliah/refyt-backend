package scheduler

import (
	"fmt"
	"refyt-backend/libs"
	"refyt-backend/scheduler/repo"
)

type Scheduler struct {
	repo repo.ISchedulerRepo
}

func NewScheduler(db *libs.PostgresDatabase) (scheduler Scheduler) {

	repo := repo.NewSchedulerRepo(db)

	return Scheduler{
		repo: repo,
	}
}

func (s *Scheduler) ProcessUpcomingBookings() (err error) {
	return nil
}

func (s *Scheduler) ProcessUpcomingPickups() (err error) {
	return nil
}

func (s *Scheduler) ProcessUpcomingReturns() (err error) {
	return nil
}

func (s *Scheduler) ProcessOverdueItems() (err error) {
	return nil
}

func (s *Scheduler) ProcessScheduledTasks() (err error) {

	err = s.ProcessUpcomingBookings()

	if err != nil {
		fmt.Println("Error processing upcoming bookings")
	}

	err = s.ProcessUpcomingPickups()

	if err != nil {
		fmt.Println("Error processing upcoming pickups")
	}

	err = s.ProcessUpcomingReturns()

	if err != nil {
		fmt.Println("Error processing upcoming returns")
	}

	err = s.ProcessOverdueItems()

	if err != nil {
		fmt.Println("Error processing overdue items")
	}

	return nil
}
