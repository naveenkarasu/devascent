public class Solution {
    public long swim_in_water(long[][] grid) {
        int n = grid.Length;
        bool[,] visited = new bool[n, n];
        var heap = new PriorityQueue<(long t, int i, int j), long>();
        heap.Enqueue((grid[0][0], 0, 0), grid[0][0]);
        visited[0, 0] = true;
        int[][] dirs = { new[]{1,0}, new[]{-1,0}, new[]{0,1}, new[]{0,-1} };
        while (heap.Count > 0) {
            var (t, i, j) = heap.Dequeue();
            if (i == n - 1 && j == n - 1) return t;
            foreach (var d in dirs) {
                int ni = i + d[0], nj = j + d[1];
                if (ni >= 0 && ni < n && nj >= 0 && nj < n && !visited[ni, nj]) {
                    visited[ni, nj] = true;
                    long nt = Math.Max(t, grid[ni][nj]);
                    heap.Enqueue((nt, ni, nj), nt);
                }
            }
        }
        return -1;
    }
}
