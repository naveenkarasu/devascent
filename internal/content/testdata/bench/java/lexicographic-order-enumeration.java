import java.util.*;

class Solution {
    public long[] lexical_order(long n) {
        List<Long> result = new ArrayList<>();
        long cur = 1;
        while (result.size() < (int) n) {
            result.add(cur);
            if (cur * 10 <= n) {
                cur *= 10;
            } else {
                while (cur % 10 == 9 || cur >= n) {
                    cur /= 10;
                }
                cur += 1;
            }
        }
        long[] arr = new long[(int) n];
        for (int i = 0; i < result.size(); i++) arr[i] = result.get(i);
        return arr;
    }
}
