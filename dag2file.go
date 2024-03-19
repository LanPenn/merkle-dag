package merkledag

// Hash2File 将哈希对应的数据从 KVStore 中读取出来，然后根据路径返回对应的文件内容。
func Hash2File(store KVStore, hash []byte, path string, hp HashPool) []byte {
	// 获取哈希函数实例
	hashFunc := hp.Get()

	// 在哈希函数中写入路径
	hashFunc.Write([]byte(path))
	// 计算路径的哈希值
	pathHash := hashFunc.Sum(nil)

	// 检查传入的哈希是否与路径的哈希相等，如果不等，返回空
	if !bytes.Equal(hash, pathHash) {
		return nil
	}

	// 从 KVStore 中检索数据
	data := store.Get(hash)

	// 返回数据（文件内容）
	return data
}
