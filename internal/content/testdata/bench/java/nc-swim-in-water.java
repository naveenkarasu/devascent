import java.util.*;

class Solution {
    public long swim_in_water(long[][] grid) {
        int n = grid.length;
        boolean[][] visited = new boolean[n][n];
        PriorityQueue<long[]> heap = new PriorityQueue<>(Comparator.comparingLong(a -> a[0]));
        heap.offer(new long[]{grid[0][0], 0, 0});
        visited[0][0] = true;
        int[][] dirs = {{1,0},{-1,0},{0,1},{0,-1}};
        while (!heap.isEmpty()) {
            long[] curr = heap.poll();
            long t = curr[0];
            int i = (int) curr[1], j = (int) curr[2];
            if (i == n - 1 && j == n - 1) return t;
            for (int[] d : dirs) {
                int ni = i + d[0], nj = j + d[1];
                if (ni >= 0 && ni < n && nj >= 0 && nj < n && !visited[ni][nj]) {
                    visited[ni][nj] = true;
                    heap.offer(new long[]{Math.max(t, grid[ni][nj]), ni, nj});
                }
            }
        }
        return -1;
    }
}
