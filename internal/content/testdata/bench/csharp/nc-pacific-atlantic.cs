using System.Collections.Generic;

public class Solution {
    public long[][] pacific_atlantic(long[][] heights) {
        if (heights == null || heights.Length == 0) return new long[0][];
        int m = heights.Length, n = heights[0].Length;

        bool[][] pac = new bool[m][];
        bool[][] atl = new bool[m][];
        for (int i = 0; i < m; i++) {
            pac[i] = new bool[n];
            atl[i] = new bool[n];
        }

        var pacQueue = new Queue<int[]>();
        var atlQueue = new Queue<int[]>();

        for (int j = 0; j < n; j++) {
            pacQueue.Enqueue(new int[]{0, j}); pac[0][j] = true;
            atlQueue.Enqueue(new int[]{m-1, j}); atl[m-1][j] = true;
        }
        for (int i = 1; i < m; i++) {
            pacQueue.Enqueue(new int[]{i, 0}); pac[i][0] = true;
        }
        for (int i = 0; i < m-1; i++) {
            atlQueue.Enqueue(new int[]{i, n-1}); atl[i][n-1] = true;
        }

        Bfs(heights, pac, pacQueue, m, n);
        Bfs(heights, atl, atlQueue, m, n);

        var result = new List<long[]>();
        for (int i = 0; i < m; i++) {
            for (int j = 0; j < n; j++) {
                if (pac[i][j] && atl[i][j]) result.Add(new long[]{i, j});
            }
        }
        return result.ToArray();
    }

    private void Bfs(long[][] heights, bool[][] visited, Queue<int[]> queue, int m, int n) {
        int[][] dirs = new int[][]{ new int[]{1,0}, new int[]{-1,0}, new int[]{0,1}, new int[]{0,-1} };
        while (queue.Count > 0) {
            int[] cur = queue.Dequeue();
            int i = cur[0], j = cur[1];
            foreach (int[] d in dirs) {
                int ni = i + d[0], nj = j + d[1];
                if (ni >= 0 && ni < m && nj >= 0 && nj < n && !visited[ni][nj]
                        && heights[ni][nj] >= heights[i][j]) {
                    visited[ni][nj] = true;
                    queue.Enqueue(new int[]{ni, nj});
                }
            }
        }
    }
}
