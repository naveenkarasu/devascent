function kth_largest(nums: number[], k: number): number {
    nums = nums.slice().sort((a, b) => b - a);
    return nums[k - 1];
}
