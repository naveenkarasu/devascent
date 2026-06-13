public class Solution {
    private long _n, _m;
    private long[,] _grid;
    private const long MOD = 1_000_000_007;

    public long count_grid_colorings(long n, long m) {
        _n = n;
        _m = m;
        _grid = new long[n, m];
        return Backtrack(0);
    }

    private long Backtrack(long pos) {
        if (pos == _n * _m) return 1;
        long r = pos / _m;
        long c = pos % _m;
        long total = 0;
        for (long color = 1; color <= 3; color++) {
            if ((r == 0 || _grid[r - 1, c] != color) &&
                (c == 0 || _grid[r, c - 1] != color)) {
                _grid[r, c] = color;
                total = (total + Backtrack(pos + 1)) % MOD;
                _grid[r, c] = 0;
            }
        }
        return total;
    }
}
