public class Solution {
    public long climb_stairs(long n) {
        long a = 1, b = 1;
        for (int i = 0; i < n; i++) {
            long tmp = a;
            a = b;
            b = tmp + b;
        }
        return a;
    }
}
