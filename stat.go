package smq

import (
	"sort"
)

func (res *Result) stat(dict *Dict) {
	// 内部数据
	res.Name = dict.Name
	res.Basic.DictLen = dict.length
	res.Words.Dist = res.wordsDist
	res.Collision.Dist = res.collDist
	res.CodeLen.Dist = res.codeDist
	// 非汉字
	tmp1 := make([]rune, 0, len(res.mapNotHan))
	for k := range res.mapNotHan {
		tmp1 = append(tmp1, k)
	}
	sort.Slice(tmp1, func(i, j int) bool {
		return tmp1[i] < tmp1[j]
	})
	res.Basic.NotHans = len(tmp1)
	res.Basic.NotHan = string(tmp1)

	// 缺字
	tmp2 := make([]rune, 0, len(res.mapLack))
	for k := range res.mapLack {
		tmp2 = append(tmp2, k)
	}
	sort.Slice(tmp2, func(i, j int) bool {
		return tmp2[i] < tmp2[j]
	})
	res.Basic.Lacks = len(tmp2)
	res.Basic.Lack = string(tmp2)
	// 上屏数
	for _, v := range res.Words.Dist {
		res.Basic.Commits += v
	}
	// 打词数
	res.Words.Commits.Count = res.Basic.Commits - res.Words.Dist[1]
	res.Words.Commits.Rate = div(res.Words.Commits.Count, res.Basic.Commits)
	res.Words.Chars.Count = res.Basic.TextLen - res.Words.Dist[1]
	res.Words.Chars.Rate = div(res.Words.Chars.Count, res.Basic.TextLen)
	// 选重数
	res.Collision.Commits.Count = res.Basic.Commits - res.Collision.Dist[1]
	res.Collision.Commits.Rate = div(res.Collision.Commits.Count, res.Basic.Commits)
	res.Collision.Chars.Rate = div(res.Collision.Chars.Count, res.Basic.TextLen)
	// 码长
	for i, v := range res.CodeLen.Dist {
		res.CodeLen.Total += i * v
	}
	res.CodeLen.PerChar = div(res.CodeLen.Total, res.Basic.TextLen)

	// keys
	for k, v := range res.mapKeys {
		key := ""
		switch k {
		case '_':
			key = "left_space"
		case '+':
			key = "right_space"
		default:
			key = string(k)
		}
		if res.Keys[key] == nil {
			res.Keys[key] = new(CaR)
		}
		res.Keys[key].Count += v
	}
	for _, v := range res.Keys {
		v.Rate = div(v.Count, res.CodeLen.Total)
	}
	// combs
	res.Combs.Equivalent = div(res.toTalEq10/10, res.Combs.Count)
	res.Combs.DoubleHit.Rate = div(res.Combs.DoubleHit.Count, res.Combs.Count)
	res.Combs.TribleHit.Rate = div(res.Combs.TribleHit.Count, res.Combs.Count)
	res.Combs.SingleSpan.Rate = div(res.Combs.SingleSpan.Count, res.Combs.Count)
	res.Combs.MultiSpan.Rate = div(res.Combs.MultiSpan.Count, res.Combs.Count)
	res.Combs.LongFingersDisturb.Rate = div(res.Combs.LongFingersDisturb.Count, res.Combs.Count)
	res.Combs.LittleFingersDisturb.Rate = div(res.Combs.LittleFingersDisturb.Count, res.Combs.Count)
	// hands
	res.Hands.Left.Count = res.Hands.LL.Count + res.Hands.LR.Count/2 + res.Hands.RL.Count/2
	res.Hands.Right.Count = res.Combs.Count - res.Hands.Left.Count
	res.Hands.Same.Count = res.Hands.LL.Count + res.Hands.RR.Count
	res.Hands.Diff.Count = res.Combs.Count - res.Hands.Same.Count
	res.Hands.Left.Rate = div(res.Hands.Left.Count, res.Combs.Count)
	res.Hands.Right.Rate = div(res.Hands.Right.Count, res.Combs.Count)
	res.Hands.Same.Rate = div(res.Hands.Same.Count, res.Combs.Count)
	res.Hands.Diff.Rate = div(res.Hands.Diff.Count, res.Combs.Count)
	// fingers
	for k, v := range res.mapKeys {
		if _, ok := keyData[k]; !ok {
			res.Fingers.Dist[10].Count += v
		} else {
			res.Fingers.Dist[keyData[k].fin].Count += v
		}
	}
	for _, v := range res.Fingers.Dist {
		if v == nil {
			continue
		}
		v.Rate = div(v.Count, res.CodeLen.Total)
	}
	res.Fingers.Same.Rate = div(res.Fingers.Same.Count, res.Combs.Count)
	res.Fingers.Diff.Count = res.Hands.Same.Count - res.Fingers.Same.Count
	res.Fingers.Diff.Rate = div(res.Fingers.Diff.Count, res.Combs.Count)
	res.Hands.LL.Rate = div(res.Hands.LL.Count, res.Combs.Count)
	res.Hands.LR.Rate = div(res.Hands.LR.Count, res.Combs.Count)
	res.Hands.RL.Rate = div(res.Hands.RL.Count, res.Combs.Count)
	res.Hands.RR.Rate = div(res.Hands.RR.Count, res.Combs.Count)
}
