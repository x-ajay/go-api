package main

import "sort"

type KthLargest struct {
	Nums []int
	K    int
}

func Constructor(k int, nums []int) KthLargest {
	sort.Ints(nums)
	return KthLargest{
		K:    k,
		Nums: nums,
	}
}

func (this *KthLargest) Add(val int) int {
	found := false
	for i := 0; i < len(this.Nums); i++ {
		if this.Nums[i] > val {
			found = true
			this.Nums = append(this.Nums[:i], append([]int{val}, this.Nums[i:]...)...)
			break
		}
	}
	if !found {
		this.Nums = append(this.Nums, val)
	}
	if len(this.Nums) < this.K {
		return 0
	}
	return this.Nums[len(this.Nums)-this.K]
}

/**
 * Your KthLargest object will be instantiated and called as such:
 * obj := Constructor(k, nums);
 * param_1 := obj.Add(val);
 */
