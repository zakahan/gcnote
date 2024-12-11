// -------------------------------------------------
// Package wrench
// Author: hanzhi
// Date: 2024/12/11
// -------------------------------------------------

package wrench

import (
	"fmt"
	"testing"
)

func TestValidString(t *testing.T) {
	x := validateIndexName("../x")
	fmt.Println(x)
	x = validateIndexName("app")
	fmt.Println(x)
	x = validateIndexName("strings")
	fmt.Println(x)
	x = validateIndexName("知识内容")
	fmt.Println(x)
}
