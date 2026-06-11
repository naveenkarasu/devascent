function find_kth_largest(nums, k) {
    // Use a min-heap of size k
    // Simple approach: sort descending and return k-th element
    const sorted = nums.slice().sort((a, b) => b - a);
    return sorted[k - 1];
}
