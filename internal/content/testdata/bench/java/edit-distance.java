import java.util.*;

class Solution {
    public long edit_distance(String str1, String str2) {
        int m = str1.length(), n = str2.length();
        long[] cur = new long[n + 1];
        for (int j = 0; j <= n; j++) cur[j] = j;
        for (int i = 1; i <= m; i++) {
            long pre = cur[0];
            cur[0] = i;
            for (int j = 1; j <= n; j++) {
                long temp = cur[j];
                if (str1.charAt(i - 1) == str2.charAt(j - 1)) {
                    cur[j] = pre;
                } else {
                    cur[j] = Math.min(pre, Math.min(cur[j - 1], cur[j])) + 1;
                }
                pre = temp;
            }
        }
        return cur[n];
    }
}
