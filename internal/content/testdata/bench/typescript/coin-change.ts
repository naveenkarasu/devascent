function coin_change(coins: number[], amount: number): number {
    const INF = amount + 1;
    const dp: number[] = new Array(amount + 1).fill(INF);
    dp[0] = 0;
    for (let a = 1; a <= amount; a++) {
        for (const c of coins) {
            if (c <= a) {
                dp[a] = Math.min(dp[a], dp[a - c] + 1);
            }
        }
    }
    return dp[amount] === INF ? -1 : dp[amount];
}
