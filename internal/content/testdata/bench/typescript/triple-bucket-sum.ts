function triple_bucket_sum(counts: number[]): number {
    const MOD = 1000000007;
    const n = counts.length;
    let ans = 0;
    for (let i = 0; i < n - 2; i++) {
        for (let j = i + 1; j < n - 1; j++) {
            for (let k = j + 1; k < n; k++) {
                ans = (ans + counts[i] * counts[j] * counts[k]) % MOD;
            }
        }
    }
    return ans;
}
