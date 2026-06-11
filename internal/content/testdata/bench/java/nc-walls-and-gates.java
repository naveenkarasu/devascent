import java.util.*;

class Solution {
    public long[][] walls_and_gates(long[][] rooms) {
        if (rooms == null || rooms.length == 0 || rooms[0].length == 0) return rooms;
        long INF = 2147483647L;
        int m = rooms.length, n = rooms[0].length;
        Queue<int[]> queue = new LinkedList<>();
        for (int i = 0; i < m; i++) {
            for (int j = 0; j < n; j++) {
                if (rooms[i][j] == 0) {
                    queue.add(new int[]{i, j});
                }
            }
        }
        int[][] dirs = {{1,0},{-1,0},{0,1},{0,-1}};
        while (!queue.isEmpty()) {
            int[] cell = queue.poll();
            int i = cell[0], j = cell[1];
            for (int[] d : dirs) {
                int ni = i + d[0], nj = j + d[1];
                if (ni >= 0 && ni < m && nj >= 0 && nj < n && rooms[ni][nj] == INF) {
                    rooms[ni][nj] = rooms[i][j] + 1;
                    queue.add(new int[]{ni, nj});
                }
            }
        }
        return rooms;
    }
}
