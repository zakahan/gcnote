// -------------------------------------------------
// Package wrench
// Author: hanzhi
// Date: 2024/12/10
// -------------------------------------------------

package wrench

import "github.com/google/uuid"

func IdGenerator() string {
	return uuid.New().String()
}
