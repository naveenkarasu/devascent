function find_kth_largest(nums: number[], k: number): number {
    nums.sort((a, b) => b - a);
    return nums[k - 1];
}
