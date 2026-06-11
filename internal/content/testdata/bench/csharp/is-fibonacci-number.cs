public class Solution {
    public bool is_fibonacci(long n) {
        if (n <= 0) return false;
        long a = 1, b = 1;
        while (a < n) {
            long temp = a + b;
            a = b;
            b = temp;
        }
        return a == n;
    }
}
