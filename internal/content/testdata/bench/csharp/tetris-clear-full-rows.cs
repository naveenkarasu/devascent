using System.Collections.Generic;
using System.Linq;

public class Solution {
    public string[][] clear_full_rows(string[][] board, string empty) {
        if (board.Length == 0) return board;
        int cols = board[0].Length;
        var surviving = new List<string[]>();
        foreach (var row in board) {
            if (row.Contains(empty)) surviving.Add(row);
        }
        int clearedCount = board.Length - surviving.Count;
        var result = new string[board.Length][];
        for (int i = 0; i < clearedCount; i++) {
            result[i] = new string[cols];
            for (int j = 0; j < cols; j++) result[i][j] = empty;
        }
        for (int i = 0; i < surviving.Count; i++) {
            result[clearedCount + i] = surviving[i];
        }
        return result;
    }
}
