package main

import (
	"fmt"
	"math/rand"
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

	//乱数
	rand.Seed(time.Now().UnixNano())

	//鍵をランダム生成
	keys := random(0b0, 0b11111111111111111111111111111111, times)

	//出力用
	//tmpOut := []int{}
	output := [32]string{}

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
				if output[i] == "" || v == output[i] {
					output[i] = v
				} else if v != output[i] {
					output[i] = "U"
				}
			}
			//fmt.Println(output)
		}
		fmt.Println(output)
		output = [32]string{}
	}

	//結果の表示
	//fmt.Printf("最終結果：%032b\n", output)

	//実行時間の表示
	fmt.Println("実行時間：", time.Since(start))
}

//chaskeyの暗号本体
func chaskey(k int, pos int, ch chan [32]string) {
	output := 0b0
	res := 0b0
	itg := [32]string{}
	var t, b, in int

	//副鍵生成
	k1 := createK1(k)
	in = rand.Intn(2)

	for i := 0b0; i <= 0b1111111111111111111111111111111; i++ { //31階差分
		//if i == 0b11111111 {
		//	fmt.Println("8階まで終了")
		//} else if i == 0b1111111111111111 {
		//	fmt.Println("16階まで終了")
		//} else if i == 0b111111111111111111111111 {
		//	fmt.Println("24階まで終了")
		//} else if i == 0b1111111111111111111111111111 {
		//	fmt.Println("28階まで終了")
		//}

		//差分ベクトルにcを差し込む処理
		t = (i >> pos) & create2(31-pos)
		b = i & create2(pos)

		if pos == 0 {
			output = ((t << 1) | in) << pos
		} else if pos == 31 {
			output = (in << 31) | b
		} else {
			output = (((t << 1) | in) << pos) | b
		}

		//鍵と平文と副鍵を排他
		output = (k ^ output) ^ k1

		//π関数
		for k := 0; k < round; k++ {
			output = permutation(output)
		}
		output ^= k1 //π関数の出力と副鍵を排他

		//1bit毎の算術和を計算
		for j := 0; j < 32; j++ {
			res ^= output
		}
	}

	//fmt.Println(res)
	for j := 0; j < 32; j++ {
		tmp := (res >> (31 - j)) & 1
		if tmp == 0b0 {
			itg[j] = "B"
		} else if tmp == 0b1 {
			itg[j] = "O"
		} else {
			fmt.Println("分岐ミス")
		}
	}
	//fmt.Println("各試行の特性", itg)
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
	res := 0b0
	switch num {
	case 0:
		res = 0b0
		break
	case 1:
		res = 0b1
		break
	case 2:
		res = 0b11
		break
	case 3:
		res = 0b111
		break
	case 4:
		res = 0b1111
		break
	case 5:
		res = 0b11111
		break
	case 6:
		res = 0b111111
		break
	case 7:
		res = 0b1111111
		break
	case 8:
		res = 0b11111111
		break
	case 9:
		res = 0b111111111
		break
	case 10:
		res = 0b1111111111
		break
	case 11:
		res = 0b11111111111
		break
	case 12:
		res = 0b111111111111
		break
	case 13:
		res = 0b1111111111111
		break
	case 14:
		res = 0b11111111111111
		break
	case 15:
		res = 0b111111111111111
		break
	case 16:
		res = 0b1111111111111111
		break
	case 17:
		res = 0b11111111111111111
		break
	case 18:
		res = 0b111111111111111111
		break
	case 19:
		res = 0b1111111111111111111
		break
	case 20:
		res = 0b11111111111111111111
		break
	case 21:
		res = 0b111111111111111111111
		break
	case 22:
		res = 0b1111111111111111111111
		break
	case 23:
		res = 0b11111111111111111111111
		break
	case 24:
		res = 0b111111111111111111111111
		break
	case 25:
		res = 0b1111111111111111111111111
		break
	case 26:
		res = 0b11111111111111111111111111
		break
	case 27:
		res = 0b111111111111111111111111111
		break
	case 28:
		res = 0b1111111111111111111111111111
		break
	case 29:
		res = 0b11111111111111111111111111111
		break
	case 30:
		res = 0b111111111111111111111111111111
		break
	case 31:
		res = 0b1111111111111111111111111111111
		break
	default:
		panic("範囲外")
	}
	return res
}
