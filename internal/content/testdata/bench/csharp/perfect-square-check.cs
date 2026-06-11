public class Solution {
    public bool is_perfect_square(long num) {
        if (num < 1) return false;
        long lo = 1, hi = num;
        while (lo <= hi) {
            long mid = (lo + hi) / 2;
            long sq = mid * mid;
            if (sq == num) return true;
            else if (sq < num) lo = mid + 1;
            else hi = mid - 1;
        }
        return false;
    }
}
