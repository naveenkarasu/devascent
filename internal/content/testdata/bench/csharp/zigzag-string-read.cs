using System.Text;

public class Solution {
    public string zigzag_convert(string s, long num_rows) {
        if (num_rows == 1) return s;
        var rows = new StringBuilder[num_rows];
        for (int i = 0; i < num_rows; i++) rows[i] = new StringBuilder();
        int row = 0, direction = -1;
        foreach (char ch in s) {
            rows[row].Append(ch);
            if (row == 0 || row == (int)num_rows - 1) {
                direction = -direction;
            }
            row += direction;
        }
        var sb = new StringBuilder();
        foreach (var r in rows) sb.Append(r);
        return sb.ToString();
    }
}
