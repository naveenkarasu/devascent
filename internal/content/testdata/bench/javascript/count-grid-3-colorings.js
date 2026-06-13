function count_grid_colorings(n, m) {
  const MOD = 1000000007;
  const grid = Array.from({length: n}, () => new Array(m).fill(0));

  function backtrack(pos) {
    if (pos === n * m) return 1;
    const r = Math.floor(pos / m);
    const c = pos % m;
    let total = 0;
    for (let color = 1; color <= 3; color++) {
      if ((r === 0 || grid[r - 1][c] !== color) &&
          (c === 0 || grid[r][c - 1] !== color)) {
        grid[r][c] = color;
        total = (total + backtrack(pos + 1)) % MOD;
        grid[r][c] = 0;
      }
    }
    return total;
  }

  return backtrack(0);
}
