import java.util.*;

class Solution {
    public long modular_exponentiation(long x, long n, long m) {
        long res = 1;
        x = x % m;
        while (n > 0) {
            if (n % 2 != 0) {
                res = (res * x) % m;
                n--;
            } else {
                x = (x * x) % m;
                n /= 2;
            }
        }
        return res;
    }
}
