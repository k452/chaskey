package main

import (
	"fmt"
	"math/rand"
	"time"
)

const blockLen = 32 //ブロック長
const keyLen = 32   //鍵長
const splitLen = 8  //
const piRound = 2   //π関数の中の転置の段数
const times = 1     //試行回数

func main() {
	//実行時間の計測開始
	start := time.Now()

	//平文をランダム生成
	texts := random(0b0, 0b11111111111111111111111111111111, times)
	keys := random(0b0, 0b11111111111111111111111111111111, times)

	//暗号化を複数回実行
	for i := 0; i < times; i++ {
		chaskey(texts[i], keys[i], i)
	}

	//実行時間の表示
	fmt.Println(time.Since(start))
}

//chaskeyの暗号本体
func chaskey(in, k, num int) {
	result := 0b00000000000000000000000000000000
	k1 := createK1(k)

	for j := 0b0; j <= 0b1111111111111111111111111111111; j++ {
		m := in ^ j      //差分ベクトルと平文を排他
		m = (k ^ m) ^ k1 //鍵と平文と副鍵を排他
		for k := 0; k < piRound; k++ {
			m = permutation(m) //π関数
		}
		m ^= k1     //π関数の出力と副鍵を排他
		result ^= m //結果をこれまでの結果と排他
	}
	fmt.Printf("試行: %d回目\n", num+1)
	fmt.Printf("結果: %032b\n", result)
	fmt.Println("----------------------------------------------------------")
	result = 0b00000000000000000000000000000000
}

//π関数の中の1ラウンドの転置
func permutation(m int) int {
	var vIn []int
	vIn = append(vIn, (m >> 24 & 0xff))
	vIn = append(vIn, (m >> 16 & 0xff))
	vIn = append(vIn, (m >> 8 & 0xff))
	vIn = append(vIn, (m & 0xff))

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

	return joinBit(joinBit(joinBit(v0_3, v1_4), v2_3), v3_4)
}

//算術和ではみ出す桁を除去
func modPlus(a int, b int) int {
	return (a + b) % 256
}

//RotateL8 8bitの数値を対象とした左シフト
func RotateL8(a int, i int) int {
	return ((a<<i)&0xff ^ (a >> (8 - i)))
}

//RotateL32 32bitの数値を対象とした左シフト
func RotateL32(a int, i int) int {
	return ((a<<i)&0xffffffff ^ (a >> (32 - i)))
}

//副鍵を生成
func createK1(k int) int {
	if ((k >> 31) & 1) == 0 {
		return RotateL32(k, 1)
	}
	return k ^ 0b00000000000000000000000010000111
}

//8bitを結合 a->何bitでもよい b->8bit
func joinBit(a int, b int) int {
	return (a << 8) | b
}

//以下2つが任意の範囲で乱数を生成して配列に格納する関数
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
	//sort.Sort(sort.IntSlice(keys)) // ソートしない場合コメントアウト
	return keys
}
