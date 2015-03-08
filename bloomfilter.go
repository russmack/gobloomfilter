package main

import (
	"fmt"
	"hash/fnv"
	"math"
)

type Hash32Fn func(string) uint32

type BloomFilter struct {
	Filter       []bool
	Size         uint32
	HashFuncs    []Hash32Fn
	TotalFlipped uint32
}

func main() {
	b := NewBloomFilter(15)
	b.Add("dog")
	t1 := "cat"
	t2 := "dog"
	r1 := b.Exists(t1)
	r2 := b.Exists(t2)
	displayResult(b, t1, r1)
	displayResult(b, t2, r2)
}

func displayResult(b *BloomFilter, s string, exists bool) {
	fmt.Println("Probability of false positives:", b.getFalsePositiveProbability())
	if !exists {
		fmt.Println("Input", s, "is certainly not in the filter.")
	} else {
		fmt.Println("Input", s, "might be in the filter.")
	}
}

// NewBloomFilter creates a new BloomFilter of the specified size.
func NewBloomFilter(filterSize uint32) *BloomFilter {
	b := BloomFilter{}
	b.Filter = make([]bool, filterSize)
	b.Size = filterSize
	b.HashFuncs = []Hash32Fn{hashFnv1, hashFnv1a}
	return &b
}

func hashFnv1(s string) uint32 {
	h := fnv.New32()
	h.Write([]byte(s))
	return h.Sum32()
}

func hashFnv1a(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// getIndex32 reduces a string hash to an index within the bounds of the filter.
func (b *BloomFilter) getIndex32(s string, hashFn Hash32Fn) uint32 {
	h := hashFn(s)
	return h % b.Size
}

// Add a string to the Bloom filter.
func (b *BloomFilter) Add(s string) {
	fmt.Println(s)
	for _, j := range b.HashFuncs {
		idx := b.getIndex32(s, j)
		fmt.Println(idx)
		b.Filter[idx] = true
		fmt.Println(b.Filter)
		b.TotalFlipped++
	}
}

// Exists checks if the given string is in the Bloom Filter.
func (b *BloomFilter) Exists(s string) bool {
	results := make([]bool, len(b.HashFuncs))
	for i, j := range b.HashFuncs {
		idx := b.getIndex32(s, j)
		val := b.Filter[idx]
		results[i] = val
	}
	allTrue := true
	trueCount := 0
	for _, j := range results {
		if j {
			trueCount++
		}
		allTrue = allTrue && j
	}
	return allTrue
}

// False positive probability:
// The number of true cells in the filter,
// divided by the length of the filter,
// to the power of the number of hash functions.
func (b *BloomFilter) getFalsePositiveProbability() float64 {
	x := float64(b.TotalFlipped) / float64(b.Size)
	y := float64(len(b.HashFuncs))
	z := math.Pow(x, y) * 100
	p := math.Floor(z + .5)
	return p
}
