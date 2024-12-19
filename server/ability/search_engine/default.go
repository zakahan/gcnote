// -------------------------------------------------
// Package search_engine
// Author: hanzhi
// Date: 2024/12/18
// -------------------------------------------------

package search_engine

type Param struct {
	Shards     int
	Replicas   int
	K1         float64
	B          float64
	Dim        int
	Similarity string
	IndexName  string
}

func Constructor() Param {
	shards := 1
	replicas := 1
	k1 := 1.2
	b := 0.75
	dim := 1024
	similarity := "dot_product"
	return Param{
		Shards:     shards,
		Replicas:   replicas,
		K1:         k1,
		B:          b,
		Dim:        dim,
		Similarity: similarity,
	}
}

var param Param = Constructor()

func _defaultSettings() map[string]interface{} {
	settings := map[string]interface{}{
		"index": map[string]interface{}{
			"number_of_shards":   param.Shards,
			"number_of_replicas": param.Replicas,
			"similarity": map[string]interface{}{
				"custom_bm25": map[string]interface{}{
					"type": "BM25",
					"k1":   param.K1,
					"b":    param.B,
				},
			},
		},
	}
	return settings
}

func _defaultMappings() map[string]interface{} {
	mappings := map[string]interface{}{
		"properties": map[string]interface{}{
			"page_content": map[string]interface{}{ // 替换为实际的常量或字符串
				"index":           true,
				"type":            "text",
				"similarity":      "custom_bm25", // 自定义 BM25 相似度
				"analyzer":        "standard",
				"search_analyzer": "standard",
			},
			"vector": map[string]interface{}{ // 替换为实际的常量或字符串
				"index":      true,
				"type":       "dense_vector",
				"dims":       param.Dim,        // 替换为实际变量
				"similarity": param.Similarity, // 替换为实际变量
			},
			"metadata": map[string]interface{}{ // 替换为实际的常量或字符串
				"type": "object",
				"properties": map[string]interface{}{
					"doc_id": map[string]interface{}{
						"type": "keyword",
					},
					"kb_file_id": map[string]interface{}{
						"type": "keyword",
					},
					"index_id": map[string]interface{}{
						"type":  "keyword",
						"index": false,
					},
					"type": map[string]interface{}{
						"type": "keyword",
					},
					"image_path": map[string]interface{}{
						"type": "keyword",
					},
				},
			},
		},
	}
	return mappings
}
