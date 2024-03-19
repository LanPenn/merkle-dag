package merkledag

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"hash"
)

// Node represents a node in the Merkle tree
type Node struct {
	Data []byte
}

// Link represents a link between nodes
type Link struct {
	Name string
	Hash []byte
	Size int
}

// Object represents an object in the Merkle DAG
type Object struct {
	Links []Link
	Data  []byte
}

// KVStore represents a key-value store
type KVStore interface {
	Has(key []byte) (bool, error)
	Put(key, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
}

// Add 函数将 Node 中的数据保存在 KVStore 中，并返回 Merkle Root
func Add(store KVStore, node Node, h hash.Hash) ([]byte, error) {
	// 计算数据的哈希值
	hashValue := h.Sum(node.Data)

	// 将数据写入 KVStore
	err := store.Put(hashValue, node.Data)
	if err != nil {
		return nil, err
	}

	// 返回 Merkle Root
	return computeMerkleRoot(store, h)
}

// computeMerkleRoot 计算 Merkle Root
func computeMerkleRoot(store KVStore, h hash.Hash) ([]byte, error) {
	var hashValues [][]byte

	// 收集 KVStore 中的所有数据的哈希值
	keys, err := getAllKeys(store)
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		data, err := store.Get(key)
		if err != nil {
			return nil, err
		}
		hashValue := h.Sum(data)
		hashValues = append(hashValues, hashValue)
	}

	// 生成 Merkle Root
	for len(hashValues) > 1 {
		var newHashValues [][]byte
		for i := 0; i < len(hashValues); i += 2 {
			concatenatedHashes := append(hashValues[i], hashValues[i+1]...)
			newHashValue := h.Sum(concatenatedHashes)
			newHashValues = append(newHashValues, newHashValue)
		}
		hashValues = newHashValues
	}

	return hashValues[0], nil
}

// getAllKeys 获取 KVStore 中的所有键
func getAllKeys(store KVStore) ([][]byte, error) {
	keys := make([][]byte, 0)
	iter := store.Iterator()
	defer iter.Release()
	for iter.Next() {
		key := iter.Key()
		keys = append(keys, key)
	}
	if err := iter.Error(); err != nil {
		return nil, err
	}
	return keys, nil
}

// bytesToHex 将字节切片转换为十六进制字符串
func bytesToHex(b []byte) string {
	var buf bytes.Buffer
	for _, v := range b {
		buf.WriteString(byteToHex(v))
	}
	return buf.String()
}

// byteToHex 将字节转换为十六进制字符串
func byteToHex(b byte) string {
	const hex = "0123456789abcdef"
	return string([]byte{hex[b>>4], hex[b&0xf]})
}
