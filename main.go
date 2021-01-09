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
const piRound = 2

var times int = int(math.Pow(2, 3)) //最終的には(2, 32)にする

func main() {
	sabun := random(0b0, 0b00000000000000000000000000000111, int(math.Pow(2, 3)))
	texts := random(0b0, 0b11111111111111111111111111111111, times)
	keys := random(0b0, 0b11111111111111111111111111111111, 1)

	for i := 0; i < times; i++ {
		result := []int{
			0b00000000,
			0b00000000,
			0b00000000,
			0b00000000,
		}
		for _, v := range sabun {
			var text []int
			text = append(text, (texts[i]>>24&0xff)^v)
			text = append(text, (texts[i]>>16&0xff)^v)
			text = append(text, (texts[i]>>8&0xff)^v)
			text = append(text, (texts[i]&0xff)^v)

			//fmt.Printf("試行: %d回目\n", i+1)
			//fmt.Printf("差分: %d個目\n", j+1)
			//fmt.Printf("平文: %08b\n", text)
			//fmt.Printf("K:    %032b\n", keys[i])
			//fmt.Printf("K1:   %032b\n", createK1(keys[i]))

			for k := 0; k < piRound; k++ {
				text = permutation(text)
			}
			//fmt.Printf("出力: %08b\n", text)
			for l := 0; l < 4; l++ {
				result[l] = text[l] ^ result[l]
			}
			//fmt.Println("----------------------------------------------------------")
		}
		fmt.Printf("試行: %d回目\n", i+1)
		fmt.Printf("結果: %08b\n", result)
		fmt.Println("----------------------------------------------------------")
		result = []int{
			0b00000000,
			0b00000000,
			0b00000000,
			0b00000000,
		}
	}
}

func permutation(vIn []int) []int {
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
	numRange := max - min + 1

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
