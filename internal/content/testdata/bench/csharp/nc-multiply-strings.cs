using System.Text;

public class Solution {
    public string multiply(string num1, string num2) {
        if (num1 == "0" || num2 == "0") return "0";
        int m = num1.Length, n = num2.Length;
        int[] buf = new int[m + n];
        for (int i = m - 1; i >= 0; i--) {
            for (int j = n - 1; j >= 0; j--) {
                int mul = (num1[i] - '0') * (num2[j] - '0');
                int p1 = i + j, p2 = i + j + 1;
                int total = mul + buf[p2];
                buf[p2] = total % 10;
                buf[p1] += total / 10;
            }
        }
        var sb = new StringBuilder();
        foreach (int d in buf) {
            if (sb.Length == 0 && d == 0) continue;
            sb.Append(d);
        }
        return sb.Length == 0 ? "0" : sb.ToString();
    }
}
