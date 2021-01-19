package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const blockLen = 32 //ブロック長
const keyLen = 32   //鍵長
const splitLen = 8  //
const round = 2     //π関数の中の転置の段数
const times = 10000 //試行回数

/* bitを1つずつ
bit := fmt.Sprintf("%04b", 0b1110)
arr := strings.Split(bit, "")
for _, v := range arr {
	num, _ := strconv.Atoi(v)
	fmt.Println(num + 10)
}
*/
func main() {
	//実行時間の計測開始
	start := time.Now()

	//差分位置の読み込み
	f, _ := os.Open("./16kai.txt")
	defer f.Close()
	scanner := bufio.NewScanner(f)

	//channelの用意
	c := make(chan [32]string)

	//出力用
	//tmpOut := []int{}
	output := [32]string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""}

	//1試行
	for scanner.Scan() {
		tmp := strings.Split(scanner.Text(), ",")
		org := []string{"0", "0", "0", "0", "0", "0", "0", "0"}
		posA := strings.Split(tmp[0], "")
		posC := strings.Split(tmp[1], "")
		fmt.Println("Aの位置", posA)

		for j := 0; j < times; j++ {
			//平文をランダム生成
			texts := random(0b0, 0b1111111111111111, times)
			keys := random(0b0, 0b11111111111111111111111111111111, times)

			text := nSplit(fmt.Sprintf("%016b", texts[j]), 4)
			go chaskey(keys[j], text, org, posA, posC, c)
		}
		for j := 0; j < times; j++ {
			for i, v := range <-c {
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
func chaskey(k int, text, org, posA, posC []string, c chan [32]string) {
	output := 0b0
	res := [32]int64{}
	itg := [32]string{}

	for i := 0b0; i <= 0b111; i++ { //16階差分
		sabun := nSplit(fmt.Sprintf("%016b", i), 4)
		for i, v := range posA {
			v, _ := strconv.Atoi(v)
			org[v] = sabun[i]
		}
		for i, v := range posC {
			v, _ := strconv.Atoi(v)
			org[v] = text[i]
		}
		pt, _ := strconv.ParseInt(strings.Join(org, ""), 2, 32)
		output = int(pt)
		k1 := createK1(k)
		output = (k ^ output) ^ k1 //鍵と平文と副鍵を排他
		for k := 0; k < round; k++ {
			output = permutation(output) //π関数
		}
		output ^= k1 //π関数の出力と副鍵を排他

		//1bit毎の算術和を計算
		for i := 0; i < 32; i++ {
			res[i] += int64((output >> (31 - i)) & 1)
		}
	}

	//fmt.Println(res)
	for q, v := range res {
		if v == 0 || v == int64(math.Pow(2, 16)) {
			itg[q] = "C"
		} else if v%2 == 0 {
			itg[q] = "B"
		} else {
			itg[q] = "U"
		}
	}
	//fmt.Println(itg)
	c <- itg
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
