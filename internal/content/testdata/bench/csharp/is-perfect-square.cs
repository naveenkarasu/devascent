public class Solution {
    public bool is_perfect_square(long n) {
        if (n < 0) return false;
        long r = (long)System.Math.Sqrt((double)n);
        while (r * r > n) r--;
        while ((r + 1) * (r + 1) <= n) r++;
        return r * r == n;
    }
}
