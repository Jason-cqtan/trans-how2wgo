// Package morestrings 实现了额外的函数来操作UTF-8
// 编码的字符串，超出了标准 "strings "包所提供的内容。
package morestrings

// ReverseRunes 返回其参数字符串从左到右颠倒的符文。
func ReverseRunes(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
