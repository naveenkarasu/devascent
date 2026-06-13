using System.Collections.Generic;

public class Solution {
    public long oranges_rotting(long[][] grid) {
        int m = grid.Length, n = grid[0].Length;
        var queue = new Queue<int[]>();
        long fresh = 0;
        for (int i = 0; i < m; i++) {
            for (int j = 0; j < n; j++) {
                if (grid[i][j] == 2) queue.Enqueue(new int[]{i, j, 0});
                else if (grid[i][j] == 1) fresh++;
            }
        }
        if (fresh == 0) return 0;
        long minutes = 0;
        int[][] dirs = new int[][]{ new int[]{1,0}, new int[]{-1,0}, new int[]{0,1}, new int[]{0,-1} };
        while (queue.Count > 0) {
            int[] cur = queue.Dequeue();
            int i = cur[0], j = cur[1], t = cur[2];
            foreach (int[] d in dirs) {
                int ni = i + d[0], nj = j + d[1];
                if (ni >= 0 && ni < m && nj >= 0 && nj < n && grid[ni][nj] == 1) {
                    grid[ni][nj] = 2;
                    fresh--;
                    minutes = t + 1;
                    queue.Enqueue(new int[]{ni, nj, t+1});
                }
            }
        }
        return fresh == 0 ? minutes : -1;
    }
}
