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
	x := ValidateIndexName("../x")
	fmt.Println(x)
	x = ValidateIndexName("app")
	fmt.Println(x)
	x = ValidateIndexName("strings")
	fmt.Println(x)
	x = ValidateIndexName("知识内容")
	fmt.Println(x)
}
