public class Solution {
    public string multiply_strings(string num1, string num2) {
        int m = num1.Length, n = num2.Length;
        int[] pos = new int[m + n];
        for (int i = m - 1; i >= 0; i--) {
            for (int j = n - 1; j >= 0; j--) {
                int mul = (num1[i] - '0') * (num2[j] - '0');
                int p1 = i + j, p2 = i + j + 1;
                int total = mul + pos[p2];
                pos[p2] = total % 10;
                pos[p1] += total / 10;
            }
        }
        var sb = new System.Text.StringBuilder();
        foreach (int d in pos) {
            if (!(sb.Length == 0 && d == 0)) sb.Append(d);
        }
        return sb.Length == 0 ? "0" : sb.ToString();
    }
}
