function num_distinct(s: string, t: string): number {
    const m = s.length;
    const n = t.length;
    const dp: number[] = new Array(n + 1).fill(0);
    dp[0] = 1;
    for (let i = 1; i <= m; i++) {
        for (let j = n; j >= 1; j--) {
            if (s[i - 1] === t[j - 1]) {
                dp[j] += dp[j - 1];
            }
        }
    }
    return dp[n];
}
