// -------------------------------------------------
// Package splitter
// Author: hanzhi
// Date: 2024/12/20
// -------------------------------------------------

package splitter

import (
	"fmt"
	"testing"
)

func TestDemo(t *testing.T) {
	x := ExtractImageURL(`![史蒂芬大教堂](/test/data/1.png "这是一个教堂")`)
	fmt.Println(x)
}
