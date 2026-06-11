function num_decodings(s) {
    const n = s.length;
    const dp = new Array(n + 1).fill(0);
    dp[0] = 1;
    dp[1] = s[0] === '0' ? 0 : 1;
    for (let i = 2; i <= n; i++) {
        const one = parseInt(s[i - 1], 10);
        const two = parseInt(s.slice(i - 2, i), 10);
        if (one !== 0) dp[i] += dp[i - 1];
        if (two >= 10 && two <= 26) dp[i] += dp[i - 2];
    }
    return dp[n];
}
