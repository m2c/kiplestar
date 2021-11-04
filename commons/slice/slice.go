package slice

func SplitStringAsChunks(str string, size int) []string {
	var chunks []string
	for i := 0; i < len(str); i += size {
		end := i + size
		if end > len(str) {
			end = len(str)
		}
		chunks = append(chunks, str[i:end])
	}
	return chunks
}
