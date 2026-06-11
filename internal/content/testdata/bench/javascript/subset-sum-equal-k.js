function subset_sum_to_k(nums, k) {
    const n = nums.length;
    let prev = new Array(k + 1).fill(false);
    prev[0] = true;
    if (nums[0] <= k) prev[nums[0]] = true;
    for (let i = 1; i < n; i++) {
        const curr = new Array(k + 1).fill(false);
        curr[0] = true;
        for (let j = 1; j <= k; j++) {
            const not_take = prev[j];
            const take = nums[i] <= j ? prev[j - nums[i]] : false;
            curr[j] = take || not_take;
        }
        prev = curr;
    }
    return prev[k];
}
