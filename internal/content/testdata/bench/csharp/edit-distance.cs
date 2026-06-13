public class Solution {
    public long edit_distance(string str1, string str2) {
        int m = str1.Length, n = str2.Length;
        long[] cur = new long[n + 1];
        for (int j = 0; j <= n; j++) cur[j] = j;
        for (int i = 1; i <= m; i++) {
            long pre = cur[0];
            cur[0] = i;
            for (int j = 1; j <= n; j++) {
                long temp = cur[j];
                if (str1[i - 1] == str2[j - 1]) {
                    cur[j] = pre;
                } else {
                    cur[j] = System.Math.Min(pre, System.Math.Min(cur[j - 1], cur[j])) + 1;
                }
                pre = temp;
            }
        }
        return cur[n];
    }
}
