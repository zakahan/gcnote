// -------------------------------------------------
// Package embeds
// Author: hanzhi
// Date: 2024/12/19
// -------------------------------------------------

package embeds

import (
	"fmt"
	"gcnote/server/ability/document"
	"math"
	"math/rand"
)

// RandEmbedding 模拟一下embedding，毕竟现在也不能来真的
func RandEmbedding(documents []*document.Document) ([][]float64, error) {
	l := len(documents)
	r := rand.New(rand.NewSource(int64(l)))
	embedding := make([][]float64, l)
	for j := range embedding {
		// 创建一个长度为1024的float64类型的切片
		randomFloats := make([]float64, 1024)
		// 填充切片，每个元素都是0到1之间的随机float64数
		for i := range randomFloats {
			randomFloats[i] = r.Float64()
		}
		randomFloats, err := NormalizeVector(randomFloats)
		if err != nil {
			return nil, err
		}
		embedding[j] = randomFloats

	}
	return embedding, nil
}

func QueryRandEmbedding(query string) ([]float64, error) {
	randomFloats := make([]float64, 1024)
	// 填充切片，每个元素都是0到1之间的随机float64数
	r := rand.New(rand.NewSource(42))
	for i := range randomFloats {
		randomFloats[i] = r.Float64()
	}
	randomFloats, err := NormalizeVector(randomFloats)
	if err != nil {
		return nil, err
	}
	return randomFloats, nil
}

func NormalizeVector(v []float64) ([]float64, error) {
	if len(v) == 0 {
		return nil, fmt.Errorf("vector cannot be empty")
	}

	// Calculate the Euclidean norm (length) of the vector.
	norm := 0.0
	for _, value := range v {
		norm += value * value
	}
	norm = math.Sqrt(norm)

	// Avoid division by zero.
	if norm == 0 {
		return nil, fmt.Errorf("vector norm cannot be zero")
	}

	// Normalize the vector.
	unitVector := make([]float64, len(v))
	for i, value := range v {
		unitVector[i] = value / norm
	}

	return unitVector, nil
}
