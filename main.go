package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

const blockLen = 32
const keyLen = 32
const splitLen = 8

var times int = int(math.Pow(2, 1))

func main() {
	texts := random(0b0, 0b11111111111111111111111111111111, times)
	//plainText := []int{0b10000000, 0b00000000, 0b00000000, 0b00000000}
	keys := random(0b0, 0b11111111111111111111111111111111, times)
	//k := 0b00000000100000001000000010000000
	//k1 := createK1(k)

	for i := 0; i < times; i++ {
		var text []int
		text = append(text, texts[i]>>24&0xff)
		text = append(text, texts[i]>>16&0xff)
		text = append(text, texts[i]>>8&0xff)
		text = append(text, texts[i]&0xff)

		fmt.Printf("乱数:   %d個目\n", i+1)
		fmt.Printf("平文: %08b\n", text)
		fmt.Printf("K:   %032b\n", keys[i])
		fmt.Printf("K1:  %032b\n", createK1(keys[i]))

		res := text
		for j := 0; j < 3; j++ {
			fmt.Printf("π関数%d段目\n", j+1)
			res = permutation(res)
		}
		fmt.Println("----------------------------------------------------------")
	}
}

func permutation(vIn []int) []int {
	//vIn = append(vIn, in>>24&0xff)
	//vIn = append(vIn, in>>16&0xff)
	//vIn = append(vIn, in>>8&0xff)
	//vIn = append(vIn, in&0xff)

	var vOut []int

	//v3
	v3_1 := RotateL8(vIn[3], 2)
	v3_2 := (modPlus(vIn[2], vIn[3]) ^ v3_1)

	//v0
	v0_1 := modPlus(vIn[0], vIn[1])
	v0_2 := RotateL8(v0_1, 4)
	v0_3 := modPlus(v3_2, v0_2)

	//v1
	v1_1 := RotateL8(vIn[1], 3)
	v1_2 := modPlus(vIn[1], vIn[0]) ^ v1_1
	v1_3 := RotateL8(v1_2, 3)
	v1_4 := modPlus(v1_2, vIn[2]) ^ v1_3

	//v2
	v2_1 := modPlus(vIn[2], vIn[3])
	v2_2 := modPlus(v1_2, v2_1)
	v2_3 := RotateL8(v2_2, 4)

	//v3
	v3_3 := RotateL8(v3_2, 1)
	v3_4 := v0_3 ^ v3_3

	vOut = append(vOut, v0_3)
	vOut = append(vOut, v1_4)
	vOut = append(vOut, v2_3)
	vOut = append(vOut, v3_4)

	fmt.Printf("%08b ", vOut[0])
	fmt.Printf("%08b ", vOut[1])
	fmt.Printf("%08b ", vOut[2])
	fmt.Printf("%08b\n", vOut[3])

	return vOut
}

func modPlus(a int, b int) int {
	return (a + b) % 256
}

func RotateL8(a int, i int) int {
	return ((a<<i)&0xff ^ (a >> (8 - i)))
}

func RotateL32(a int, i int) int {
	return ((a<<i)&0xffffffff ^ (a >> (32 - i)))
}

func createK1(k int) int {
	if ((k >> 31) & 1) == 0 {
		return RotateL32(k, 1)
	} else {
		return k ^ 0b00000000000000000000000010000111
	}
}

func createK2(k1 int) int {
	return createK1(createK1(k1))
}

func k(m map[int]bool) []int {
	i := 0
	result := make([]int, len(m))
	for key, _ := range m {
		result[i] = key
		i++
	}
	return result
}

func random(min int, max int, num int) []int {
	numRange := max - min

	selected := make(map[int]bool)
	rand.Seed(time.Now().UnixNano())
	for counter := 0; counter < num; {
		n := rand.Intn(numRange) + min
		if selected[n] == false {
			selected[n] = true
			counter++
		}
	}
	keys := k(selected)
	// ソートしたくない場合は以下1行をコメントアウト
	sort.Sort(sort.IntSlice(keys))
	return keys
}

/*
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

//StoI string型をint型にする処理を関数化
func StoI(s string) int {
	res, _ := strconv.Atoi(s)
	return res
}

//ItoS int型をstring型にする処理を関数化
func ItoS(i int) string {
	return strconv.Itoa(i)
}

func _2to10(s string) int {
	res, _ := strconv.ParseInt(s, 2, 0)
	return int(res)
}

func _10to2(i int) string {
	return fmt.Sprintf("%b", i)
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

func rotateR(a []string, i int) string {
	return rotateL(a, len(a)-i)
}
*/
