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

func TestDemo2(t *testing.T) {
	chunks := []string{
		"![image1.png](images\\image1.png)",
		"![image2.png](images\\image2.png)",
		"![image3.png](images\\image3.png)",
	}
	x := ChunkRead(
		chunks,
		"http://127.0.0.1:8086/images",
		"514b6721-26f9-46da-9be1-0b92261d2290",
		"d6d6206b-f200-4715-86b4-b8a512fc401e",
	)
	fmt.Println(x)
}

func TestDemo3(t *testing.T) {
	chunks := []string{
		//"![image1.png](http://127.0.0.1:8086/images/514b6721-26f9-46da-9be1-0b92261d2290/d6d6206b-f200-4715-86b4-b8a512fc401e/image1.png)",
		"![image2.png](http://127.0.0.1:8086/images///1735550182_image.png)",
	}
	x, _ := ChunkReadReverse(chunks,
		"http://127.0.0.1:8086/images",
		"ccc",
		"ccc",
	)
	fmt.Println(x)
}
