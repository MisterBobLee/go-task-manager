package models

import (
	"fmt"
	"math/rand"
)

func SeedTasks(userID uint, total int) {
    var tasks []Task

    for i := 0; i < total; i++ {
        tasks = append(tasks, Task{
            Title:     fmt.Sprintf("Task #%d", i+1),
            Completed: rand.Intn(2) == 1,
            UserID:    userID,
        })
    }

    DB.CreateInBatches(tasks, 1000) // 每次批次寫入 1000 筆
}
