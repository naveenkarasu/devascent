function find_target_sum_ways(nums: number[], target: number): number {
    let dp = new Map<number, number>();
    dp.set(0, 1);
    for (const num of nums) {
        const next_dp = new Map<number, number>();
        for (const [s, cnt] of dp) {
            next_dp.set(s + num, (next_dp.get(s + num) ?? 0) + cnt);
            next_dp.set(s - num, (next_dp.get(s - num) ?? 0) + cnt);
        }
        dp = next_dp;
    }
    return dp.get(target) ?? 0;
}
