package helper

import "unicode/utf16"

// UTF16Len は文字列のUTF-16コードユニット数を返す。
// HTMLのinput要素の文字数カウントと同等。
// BMP外の文字（一部の絵文字など）はサロゲートペアで2としてカウントされる。
func UTF16Len(s string) int {
	return len(utf16.Encode([]rune(s)))
}
