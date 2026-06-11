function kth_largest(nums, k) {
    nums = nums.slice().sort((a, b) => b - a);
    return nums[k - 1];
}
