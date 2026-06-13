public class Solution {
    public long longest_ones_run(long n) {
        string binary = Convert.ToString(n, 2);
        string[] runs = binary.Split('0');
        long max = 0;
        foreach (var r in runs) {
            if (r.Length > max) max = r.Length;
        }
        return max;
    }
}
