function find_kth_largest(nums, k) {
    nums.sort((a, b) => b - a);
    return nums[k - 1];
}
