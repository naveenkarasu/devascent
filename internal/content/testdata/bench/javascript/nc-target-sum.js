function find_target_sum_ways(nums, target) {
    let dp = new Map([[0, 1]]);
    for (const num of nums) {
        const next_dp = new Map();
        for (const [s, cnt] of dp) {
            const a = s + num;
            const b = s - num;
            next_dp.set(a, (next_dp.get(a) || 0) + cnt);
            next_dp.set(b, (next_dp.get(b) || 0) + cnt);
        }
        dp = next_dp;
    }
    return dp.get(target) || 0;
}
