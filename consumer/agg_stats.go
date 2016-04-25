package main

import (
	"github.com/petar/GoLLRB/llrb"
)

type Hit struct {
	Key     string
	Counter int
}

func (h *Hit) Less(other llrb.Item) bool {
	return h.Counter < other.(*Hit).Counter
}

type HitCounter struct {
	hitMap map[string]*Hit
}

func NewHitCounter() *HitCounter {
	ret := &HitCounter{
		hitMap: make(map[string]*Hit),
	}

	return ret
}

func (hc *HitCounter) Push(key string) {
	hit, ok := hc.hitMap[key]
	if !ok {
		hit = &Hit{Key: key, Counter: 0}
		hc.hitMap[key] = hit
	}
	hit.Counter += 1
}

func (hc *HitCounter) TopN(n int) []*Hit {
	tree := llrb.New()
	for _, item := range hc.hitMap {
		tree.InsertNoReplace(item)
	}
	ret := []*Hit{}
	t := 0

	topHit, ok := tree.Max().(*Hit)
	if !ok {
		return ret
	}

	tree.DescendLessOrEqual(topHit, func(i llrb.Item) bool {
		hit := i.(*Hit)
		ret = append(ret, hit)
		t += 1
		if t >= n {
			return false
		}
		return true
	})

	return ret
}
