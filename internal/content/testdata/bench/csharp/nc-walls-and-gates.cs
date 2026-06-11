using System.Collections.Generic;

public class Solution {
    public long[][] walls_and_gates(long[][] rooms) {
        if (rooms == null || rooms.Length == 0 || rooms[0].Length == 0) return rooms;
        long INF = 2147483647;
        int m = rooms.Length, n = rooms[0].Length;
        var queue = new Queue<(int, int)>();
        for (int i = 0; i < m; i++)
            for (int j = 0; j < n; j++)
                if (rooms[i][j] == 0) queue.Enqueue((i, j));

        int[] dr = { 1, -1, 0, 0 };
        int[] dc = { 0, 0, 1, -1 };
        while (queue.Count > 0) {
            var (i, j) = queue.Dequeue();
            for (int d = 0; d < 4; d++) {
                int ni = i + dr[d], nj = j + dc[d];
                if (ni >= 0 && ni < m && nj >= 0 && nj < n && rooms[ni][nj] == INF) {
                    rooms[ni][nj] = rooms[i][j] + 1;
                    queue.Enqueue((ni, nj));
                }
            }
        }
        return rooms;
    }
}
