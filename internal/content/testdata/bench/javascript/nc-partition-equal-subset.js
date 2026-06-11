function can_partition(nums) {
    const total = nums.reduce((a, b) => a + b, 0);
    if (total % 2 !== 0) return false;
    const target = total / 2;
    let dp = new Set([0]);
    for (const n of nums) {
        const next_dp = new Set(dp);
        for (const s of dp) {
            next_dp.add(s + n);
        }
        dp = next_dp;
        if (dp.has(target)) return true;
    }
    return false;
}
