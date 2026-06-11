import java.util.*;

class Solution {
    public boolean is_fibonacci(long n) {
        if (n <= 0) return false;
        long a = 1, b = 1;
        while (a < n) {
            long tmp = a;
            a = b;
            b = tmp + b;
        }
        return a == n;
    }
}
