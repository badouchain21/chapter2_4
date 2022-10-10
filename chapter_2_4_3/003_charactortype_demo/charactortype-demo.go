package main

import "fmt"

func main() {
	Charactortype("werhewhrw232323 343434")
}
/**
输入一行字符,分别统计出其中英文字母、空格、数字和其它字符的个数。
*/
func Charactortype(str string) {

	var e,s,d,o int
	for i := o; i < len(str); i++ {
		switch {
		case 64 < str[i] && str[i] < 91: // 英文
			e += 1
		case 96 < str[i] && str[i] < 123: // 英文
			e += 1
		case 47 < str[i] && str[i] < 58: //数字
			d += 1
		case str[i] == 32: //空格
			s += 1
		default:
			o += 1
		}
	}
	fmt.Printf("字符串英文字符个数是: %d\n",e)
	fmt.Printf("字符串数字字符个数是: %d\n",d)
	fmt.Printf("字符串空格字符个数是: %d\n",s)
	fmt.Printf("字符串其它字符个数是: %d\n",o)
}