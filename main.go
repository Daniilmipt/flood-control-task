package main

import (
	"context"
	"fmt"
	"sync"
	"task/Flood"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	// create a new flood control instance
	flood := Flood.NewFlood(2, 15)

	// There is example, where we run 3 gorutins with userId = 1.
	// First and second gorutine will return true as result of Check.
	// Third will return false, due to flood have only 2 max size.
	// 4th gorutine will return true, because we didn't touch userId = 2.
	for i := 0; i < 4; i++ {
		time.Sleep(2 * time.Second)

		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			// set the user id based on the iteration number
			userId := int64(1)
			if i == 3 {
				userId = 2
			}

			ch, err := flood.Check(ctx, userId)

			fmt.Println("=========")
			fmt.Println("result for user: ", userId, ch)
			fmt.Println("error for user: ", userId, err)
			fmt.Print("--------\n\n")
		}(i)
	}

	// wait for all goroutines to finish.
	wg.Wait()

	// wait 15 seconds before checking the status for user 1 again.
	time.Sleep(15 * time.Second)

	// check if user 1 is allowed to perform a task.
	// Will return true, because it's been 15 seconds
	// and he can't find the function calls in the last 15 seconds.
	ch, err := flood.Check(ctx, 1)

	fmt.Println("=========")
	fmt.Println("result for user: ", 1, ch)
	fmt.Println("error for user: ", 1, err)
	fmt.Println("--------")
}

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}
