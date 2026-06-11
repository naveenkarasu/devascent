import java.util.*;

class Solution {
    public boolean is_valid_sudoku(String[][] board) {
        Set<String>[] rows = new HashSet[9];
        Set<String>[] cols = new HashSet[9];
        Set<String>[] boxes = new HashSet[9];
        for (int i = 0; i < 9; i++) {
            rows[i] = new HashSet<>();
            cols[i] = new HashSet<>();
            boxes[i] = new HashSet<>();
        }
        for (int r = 0; r < 9; r++) {
            for (int c = 0; c < 9; c++) {
                String val = board[r][c];
                if (val.equals(".")) continue;
                int boxIdx = (r / 3) * 3 + (c / 3);
                if (rows[r].contains(val) || cols[c].contains(val) || boxes[boxIdx].contains(val)) {
                    return false;
                }
                rows[r].add(val);
                cols[c].add(val);
                boxes[boxIdx].add(val);
            }
        }
        return true;
    }
}
