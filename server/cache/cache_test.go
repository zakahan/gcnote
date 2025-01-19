// -------------------------------------------------
// Package cache
// Author: hanzhi
// Date: 2025/1/19
// -------------------------------------------------

package cache

import (
	"context"
	"fmt"
	"gcnote/server"
	"testing"
)

func TestDequeue(t *testing.T) {
	server.InitConfig()
	server.InitRedis()

	var ctx = context.Background()
	tasks, err := DequeueAllTasks(ctx, "49c3f28c-57fc-4625-9e0c-8bf48ee0913c")
	if err != nil {
		fmt.Println("任务失败")
		return
	}
	fmt.Println(tasks)
}
