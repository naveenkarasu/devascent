import java.util.*;

class Solution {
    public long[][] pacific_atlantic(long[][] heights) {
        if (heights == null || heights.length == 0) return new long[0][0];
        int m = heights.length, n = heights[0].length;

        boolean[][] pac = new boolean[m][n];
        boolean[][] atl = new boolean[m][n];

        Queue<int[]> pacQueue = new LinkedList<>();
        Queue<int[]> atlQueue = new LinkedList<>();

        for (int j = 0; j < n; j++) {
            pacQueue.add(new int[]{0, j});
            pac[0][j] = true;
            atlQueue.add(new int[]{m-1, j});
            atl[m-1][j] = true;
        }
        for (int i = 1; i < m; i++) {
            pacQueue.add(new int[]{i, 0});
            pac[i][0] = true;
        }
        for (int i = 0; i < m-1; i++) {
            atlQueue.add(new int[]{i, n-1});
            atl[i][n-1] = true;
        }

        bfs(heights, pac, pacQueue, m, n);
        bfs(heights, atl, atlQueue, m, n);

        List<long[]> result = new ArrayList<>();
        for (int i = 0; i < m; i++) {
            for (int j = 0; j < n; j++) {
                if (pac[i][j] && atl[i][j]) {
                    result.add(new long[]{i, j});
                }
            }
        }
        return result.toArray(new long[0][]);
    }

    private void bfs(long[][] heights, boolean[][] visited, Queue<int[]> queue, int m, int n) {
        int[][] dirs = {{1,0},{-1,0},{0,1},{0,-1}};
        while (!queue.isEmpty()) {
            int[] cur = queue.poll();
            int i = cur[0], j = cur[1];
            for (int[] d : dirs) {
                int ni = i + d[0], nj = j + d[1];
                if (ni >= 0 && ni < m && nj >= 0 && nj < n && !visited[ni][nj]
                        && heights[ni][nj] >= heights[i][j]) {
                    visited[ni][nj] = true;
                    queue.add(new int[]{ni, nj});
                }
            }
        }
    }
}
