package merkledag

import (
	"bytes"
	"encoding/json"
	"strings"
)

const STEP = 4

// Hash2File 将哈希对应的数据从 KVStore 中读取出来，然后根据路径返回对应的文件内容。
func Hash2File(store KVStore, hash []byte, path string, hp HashPool) []byte {
	// 获取哈希函数实例
	hashFunc := hp.Get()

	// 在哈希函数中写入路径
	hashFunc.Write([]byte(path))
	// 计算路径的哈希值
	pathHash := hashFunc.Sum(nil)

	// 检查传入的哈希是否与路径的哈希相等，如果相等，则直接返回哈希对应的数据
	if bytes.Equal(hash, pathHash) {
		data := store.Get(hash)
		return data
	}

	// 检查路径是否合法
	if !isValidPath(path) {
		return nil
	}

	// 根据hash和path，返回对应的文件内容，hash对应的类型是tree
	return getFileFromTree(store, hash, path)
}

// 检查路径是否合法
func isValidPath(path string) bool {
	// 在此处添加任何路径验证逻辑，例如检查路径是否为空或格式是否正确
	// 这里仅简单地检查路径是否为空
	return path != ""
}

// 根据路径从树结构中获取文件内容
func getFileFromTree(store KVStore, hash []byte, path string) []byte {
	// 检索树结构对象
	treeObjBinary, _ := store.Get(hash)
	treeObj := binaryToObj(treeObjBinary)

	// 将路径拆分为路径段
	pathSegments := strings.Split(path, "/")

	// 递归查找文件
	return getFileByDir(store, treeObj, pathSegments, 1)
}

// 递归查找文件
func getFileByDir(store KVStore, obj *Object, pathSegments []string, curIndex int) []byte {
	if curIndex >= len(pathSegments) {
		return nil
	}

	// 遍历对象的链接
	for _, link := range obj.Links {
		// 检查链接名称是否匹配当前路径段
		if link.Name == pathSegments[curIndex] {
			// 根据链接类型进行处理
			switch link.Type {
			case FILE:
				// 如果是文件类型，则直接返回文件内容
				return store.Get(link.Hash)
			case DIR:
				// 如果是目录类型，则继续递归查找
				dirObjBinary, _ := store.Get(link.Hash)
				dirObj := binaryToObj(dirObjBinary)
				return getFileByDir(store, dirObj, pathSegments, curIndex+1)
			}
		}
	}

	// 未找到匹配的链接
	return nil
}

// 将二进制数据解析为对象
func binaryToObj(objBinary []byte) *Object {
	var res Object
	json.Unmarshal(objBinary, &res)
	return &res
}
