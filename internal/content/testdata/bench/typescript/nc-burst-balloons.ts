function max_coins(nums: number[]): number {
    const arr = [1, ...nums, 1];
    const n = arr.length;
    const dp: number[][] = Array.from({length: n}, () => new Array(n).fill(0));
    for (let length = 2; length < n; length++) {
        for (let left = 0; left < n - length; left++) {
            const right = left + length;
            for (let k = left + 1; k < right; k++) {
                const coins = arr[left] * arr[k] * arr[right] + dp[left][k] + dp[k][right];
                if (coins > dp[left][right]) dp[left][right] = coins;
            }
        }
    }
    return dp[0][n - 1];
}
