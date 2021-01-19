package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

const blockLen = 32 //ブロック長
const keyLen = 32   //鍵長
const splitLen = 8  //
const round = 8     //π関数の中の転置の段数
const times = 10    //試行回数

func main() {
	//実行時間の計測開始
	start := time.Now()

	//channelの用意
	ch := make(chan [32]string)

	//鍵をランダム生成
	keys := random(0b0, 0b11111111111111111111111111111111, times)

	//出力用
	//tmpOut := []int{}
	output := [32]string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""}

	//全体
	for c := 0; c < 32; c++ {
		fmt.Println("cの位置", 31-c)

		//timesの分だけ試行
		for j := 0; j < times; j++ {
			go chaskey(keys[j], c, ch)
		}

		//並列で実行した結果を最終結果としてまとめる
		for j := 0; j < times; j++ {
			for i, v := range <-ch {
				if output[i] == "" {
					output[i] = v
				} else if v == output[i] && v == "B" {
					output[i] = "B"
				} else if v == output[i] && v == "C" {
					output[i] = "C"
				} else {
					output[i] = "U"
				}
			}
			//fmt.Println(output)
		}
		fmt.Println(output)
		output = [32]string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""}
	}

	//結果の表示
	//fmt.Printf("最終結果：%032b\n", output)

	//実行時間の表示
	fmt.Println("実行時間：", time.Since(start))
}

//chaskeyの暗号本体
func chaskey(k int, pos int, ch chan [32]string) {
	output := 0b0
	res := [32]int64{}
	itg := [32]string{}

	for i := 0b0; i <= 0b1111111111111111111111111111111; i++ { //31階差分
		//差分ベクトルにcを差し込む処理
		t := (i >> pos) & create2(31-pos)
		b := i & create2(pos)
		rand.Seed(time.Now().UnixNano())
		in := rand.Intn(2)
		output = (((t << 1) | in) << i) | b

		//副鍵生成
		k1 := createK1(k)

		//鍵と平文と副鍵を排他
		output = (k ^ output) ^ k1

		//π関数
		for k := 0; k < round; k++ {
			output = permutation(output)
		}
		output ^= k1 //π関数の出力と副鍵を排他

		//1bit毎の算術和を計算
		for i := 0; i < 32; i++ {
			res[i] += int64((output >> (31 - i)) & 1)
		}
	}

	//fmt.Println(res)
	for q, v := range res {
		if v == 0 || v == int64(math.Pow(2, 31)-1) {
			itg[q] = "C"
		} else if v%2 == 0 {
			itg[q] = "B"
		} else {
			itg[q] = "U"
		}
	}
	//fmt.Println(itg)
	ch <- itg
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

//01の乱数生成
func binaryRand() string {
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(2) == 0 {
		return "0"
	}
	return "1"
}

//以下2つが任意の範囲で乱数を生成して配列に格納する関数
func k(m map[int]bool) []int {
	i := 0
	result := make([]int, len(m))
	for key := range m {
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

//文字列をn分割する
func nSplit(msg string, splitlen int) []string {
	slc := []string{}
	for i := 0; i < len(msg); i += splitlen {
		if i+splitlen < len(msg) {
			slc = append(slc, msg[i:(i+splitlen)])
		} else {
			slc = append(slc, msg[i:])
		}
	}
	return slc
}

//任意の長さの2進数(全部1)を返す
func create2(num int) int {
	txt := ""
	for j := 0; j < num; j++ {
		txt += "1"
	}
	res, _ := strconv.ParseInt(txt, 2, 32)
	return int(res)
}
