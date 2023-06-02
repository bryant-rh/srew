package util

import (
	"strings"
)

/*
	简易版的版本号对比，要求必须版本号位数相同，否则对比不了，这里也是存在一个bug，这个版本已经解决
	1.1.1 --> 1.2.1 ok
	1.1.1 --> 1.2.12 ok
	1.1.1 -->  1.2   ok
	1.2 --> 3.2.2 ok
*/

var (
	version0 = 0 //版本相等
	version1 = 1 //v1 > v2
	version2 = 2 //v1 < v2
)

func StrTrimSpace(v1str, v2str string) (v1, v2 string) {
	v1 = strings.TrimSpace(v1str)
	v2 = strings.TrimSpace(v2str)
	return
}
func comparSlice(v1slice, v2slice []string) int {
	for index := range v1slice {
		if v1slice[index] > v2slice[index] {
			return version1
		}
		if v1slice[index] < v2slice[index] {
			return version2
		}
		if len(v1slice)-1 == index {
			return version0
		}
	}
	return version0
}

func comparSlice1(v1slice, v2slice []string, flas int) int {
	for index := range v1slice {
		//按照正常逻辑v1slice 长度小
		if v1slice[index] > v2slice[index] {
			if flas == 2 {
				return version2
			}
			return version1

		}
		if v1slice[index] < v2slice[index] {
			if flas == 2 {
				return version1
			}
			return version2
		}
		if len(v1slice)-1 == index {
			if flas == 2 {
				return version1
			} else if flas == 1 {
				return version2
			}
		}
	}
	return version0
}

func CompareStrVer(v1, v2 string) (res int) {
	s1, s2 := StrTrimSpace(v1, v2)
	v1slice := strings.Split(s1, ".")
	v2slice := strings.Split(s2, ".")
	//长度不相等直接退出
	if len(v1slice) != len(v2slice) {
		if len(v1slice) > len(v2slice) {
			res = comparSlice1(v2slice, v1slice, 2)
			return res
		} else {
			res = comparSlice1(v1slice, v2slice, 1)
			return res
		}
	} else {
		res = comparSlice(v1slice, v2slice)
	}
	return res

}

// func demo01() {
// 	v1 := "5.4.0.0"
// 	v2 := "5.4.0.0"
// 	fmt.Println(compareStrVer(v1, v2))
// }
