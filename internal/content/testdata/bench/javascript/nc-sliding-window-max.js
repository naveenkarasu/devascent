function max_sliding_window(nums, k) {
    const n = nums.length;
    const result = [];
    const dq = []; // indices, decreasing by value
    for (let i = 0; i < n; i++) {
        while (dq.length > 0 && dq[0] < i - k + 1) {
            dq.shift();
        }
        while (dq.length > 0 && nums[dq[dq.length - 1]] < nums[i]) {
            dq.pop();
        }
        dq.push(i);
        if (i >= k - 1) {
            result.push(nums[dq[0]]);
        }
    }
    return result;
}
