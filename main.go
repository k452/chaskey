package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

const blockLen = 32
const keyLen = 32
const splitLen = 8

func main() {
	plainText := "00101000100000000010100000111110"
	fmt.Println(permutation(plainText))
}

func permutation(in string) []string {
	var v []string
	var vResult []string
	if utf8.RuneCountInString(in)%splitLen == 0 && utf8.RuneCountInString(in) == blockLen {
		msg := in
		runes := []rune(msg)
		for i := 0; i < len(runes); i += splitLen {
			if i+splitLen < len(runes) {
				mTmp := string(runes[i:(i + splitLen)])
				v = append(v, mTmp)
			} else {
				mTmp := string(runes[i:])
				v = append(v, mTmp)
			}
		}
	} else {
		fmt.Println("桁数：", utf8.RuneCountInString(in))
		fmt.Println(("桁数≠ブロック長"))
	}

	//v3
	v3_1 := rotateL(strings.Split(v[3], ""), 8)
	v3_2 := _10to2((_2to10(v[2]) + _2to10(v[3])) ^ StoI(v3_1))

	//v0
	v0_1 := _10to2(_2to10(v[0]) + _2to10(v[1]))
	v0_2 := rotateL(strings.Split(v0_1, ""), 16)
	v0_3 := _10to2(_2to10((v3_2)) + _2to10(v0_2))

	//v1
	v1_1 := rotateL(strings.Split(v[1], ""), 15)
	v1_2 := _10to2((_2to10(v[1]) + _2to10(v[0])) ^ StoI(v1_1))
	v1_3 := rotateL(strings.Split(v1_2, ""), 7)
	v1_4 := _10to2((_2to10(v1_2) + _2to10(v[2])) ^ _2to10(v1_3))

	//v2
	v2_1 := _10to2(_2to10(v[2]) + _2to10(v[3]))
	v2_2 := _10to2(_2to10((v1_2)) + _2to10(v2_1))
	v2_3 := rotateL(strings.Split(v2_2, ""), 16)

	//v3
	v3_3 := rotateL(strings.Split(v3_2, ""), 13)
	v3_4 := _10to2(_2to10(v0_3) ^ _2to10(v3_3))

	vResult = append(vResult, v0_3)
	vResult = append(vResult, v1_4)
	vResult = append(vResult, v2_3)
	vResult = append(vResult, v3_4)
	return vResult
}

//StoI string型をInt64型にする処理を関数化
func StoI(s string) int64 {
	res, _ := strconv.Atoi(s)
	return int64(res)
}

//ItoS int型をstring型にする処理を関数化
func ItoS(i int) string {
	return strconv.Itoa(i)
}

func _2to10(s string) int64 {
	res, _ := strconv.ParseInt(s, 2, 0)
	return res
}

func _10to2(i int64) string {
	return fmt.Sprintf("%b", i)
}

func joinArray(ary []string) string {
	s := ""
	for _, v := range ary {
		s += v
	}
	return s
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func rotateL(a []string, i int) string {
	i = i % len(a)
	if i < 0 {
		i += len(a)
	}

	for c := 0; c < gcd(i, len(a)); c++ {
		t := a[c]
		j := c
		for {
			k := j + i
			if k >= len(a) {
				k -= len(a)
			}
			if k == c {
				break
			}
			a[j] = a[k]
			j = k
		}
		a[j] = t
	}
	return joinArray((a))
}

/*
func rotateR(a []string, i int) string {
	return rotateL(a, len(a)-i)
}
*/

func createK1() {

}

func createK2() {

}
