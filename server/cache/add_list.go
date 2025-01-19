// -------------------------------------------------
// Package cache
// Author: hanzhi
// Date: 2025/1/19
// -------------------------------------------------

package cache

import (
	"context"
	"fmt"
	"gcnote/server/config"
	"log"
	"time"
)

type Task struct {
	KbFileName     string `json:"kbFileName"`
	TaskCreateTime string `json:"taskCreateTime"`
	State          string `json:"state"`
	Reason         string `json:"reason"`
}

// EnqueueTask adds a new task to the specified user's queue.
func EnqueueTask(ctx context.Context, userID string, task Task) (string, error) {
	// Generate unique task ID
	taskID := fmt.Sprintf("task:%s:%d", userID, time.Now().UnixNano())
	rdb := config.RedisClient
	// Add task details as a hash with the unique task ID
	err := rdb.HSet(ctx, taskID, map[string]interface{}{
		"kbFileName":     task.KbFileName,
		"TaskCreateTime": task.TaskCreateTime,
		"State":          task.State,
		"Reason":         task.Reason,
	}).Err()
	if err != nil {
		return "", err
	}

	// Push the task ID into the user's list(queue)
	err = rdb.LPush(ctx, userID, taskID).Err()
	if err != nil {
		return "", err
	}

	return taskID, nil
}

// DequeueAllTasks fetches all tasks from the specified user's queue and returns them.
func DequeueAllTasks(ctx context.Context, userID string) ([]map[string]string, error) {
	var tasks []map[string]string
	rdb := config.RedisClient
	// Get all task IDs from the user's queue
	taskIDs, err := rdb.LRange(ctx, userID, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	for _, taskID := range taskIDs {
		// Fetch the task details by its ID
		vals, err := rdb.HGetAll(ctx, taskID).Result()
		if err != nil {
			log.Printf("Failed to get task details for taskID %s: %v", taskID, err)
			continue
		}
		tasks = append(tasks, vals)
	}

	return tasks, nil
}
