using System.Collections.Generic;

public class Solution {
    public bool is_valid_sudoku(string[][] board) {
        var rows = new HashSet<string>[9];
        var cols = new HashSet<string>[9];
        var boxes = new HashSet<string>[9];
        for (int i = 0; i < 9; i++) {
            rows[i] = new HashSet<string>();
            cols[i] = new HashSet<string>();
            boxes[i] = new HashSet<string>();
        }
        for (int r = 0; r < 9; r++) {
            for (int c = 0; c < 9; c++) {
                string val = board[r][c];
                if (val == ".") continue;
                int boxIdx = (r / 3) * 3 + (c / 3);
                if (!rows[r].Add(val) || !cols[c].Add(val) || !boxes[boxIdx].Add(val))
                    return false;
            }
        }
        return true;
    }
}
