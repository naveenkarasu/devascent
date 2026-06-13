import java.util.*;

class Solution {
    public String[][] clear_full_rows(String[][] board, String empty) {
        if (board.length == 0) return board;
        int cols = board[0].length;
        List<String[]> surviving = new ArrayList<>();
        for (String[] row : board) {
            boolean hasEmpty = false;
            for (String cell : row) {
                if (cell.equals(empty)) { hasEmpty = true; break; }
            }
            if (hasEmpty) surviving.add(row);
        }
        int clearedCount = board.length - surviving.size();
        String[][] result = new String[board.length][cols];
        for (int i = 0; i < clearedCount; i++) {
            Arrays.fill(result[i], empty);
        }
        for (int i = 0; i < surviving.size(); i++) {
            result[clearedCount + i] = surviving.get(i);
        }
        return result;
    }
}
