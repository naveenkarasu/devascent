function min_path_sum(grid) {
  const m = grid.length, n = grid[0].length;
  const dp = grid.map(row => row.slice());
  for (let i = 0; i < m; i++) {
    for (let j = 0; j < n; j++) {
      if (i === 0 && j === 0) continue;
      else if (i === 0) dp[i][j] += dp[i][j - 1];
      else if (j === 0) dp[i][j] += dp[i - 1][j];
      else dp[i][j] += Math.min(dp[i - 1][j], dp[i][j - 1]);
    }
  }
  return dp[m - 1][n - 1];
}
