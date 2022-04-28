package smq

type key struct {
	key byte
	lor bool // Left or Right
	fin int  // 0-10, 10是错误的，左右大拇指_+
}

var keyData = newKeyData()

func newKeyData() map[byte]key {
	ret := make(map[byte]key)
	// 左手空格"_" 右手空格"+"
	lk := "_12345qwertasdfgzxcvb"  // left keys
	rk := "+67890yuiophjkl;'nm,./" // right keys
	for i := range lk {
		ret[lk[i]] = key{lk[i], false, 0}
	}
	for i := range rk {
		ret[rk[i]] = key{rk[i], true, 0}
	}
	// 手指
	keys := "1qaz2wsx3edc4rfv5tgb_+6yhn7ujm8ik,9ol.0p;/'"
	fin := "1111222233334444444456777777778888999900000"
	for i := range keys {
		tmp := ret[keys[i]]
		tmp.fin = int(fin[i] - 48)
		ret[keys[i]] = tmp
	}
	return ret
}
