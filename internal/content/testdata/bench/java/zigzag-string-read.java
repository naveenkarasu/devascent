class Solution {
    public String zigzag_convert(String s, long num_rows) {
        int numRows = (int) num_rows;
        if (numRows == 1) return s;
        StringBuilder[] rows = new StringBuilder[numRows];
        for (int i = 0; i < numRows; i++) rows[i] = new StringBuilder();
        int row = 0, direction = -1;
        for (char ch : s.toCharArray()) {
            rows[row].append(ch);
            if (row == 0 || row == numRows - 1) direction = -direction;
            row += direction;
        }
        StringBuilder result = new StringBuilder();
        for (StringBuilder r : rows) result.append(r);
        return result.toString();
    }
}
