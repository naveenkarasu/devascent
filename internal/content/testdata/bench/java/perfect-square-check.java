class Solution {
    public boolean is_perfect_square(long num) {
        if (num < 1) return false;
        long lo = 1, hi = num;
        while (lo <= hi) {
            long mid = lo + (hi - lo) / 2;
            if (mid == num / mid && num % mid == 0 && mid * mid == num) return true;
            else if (mid < num / mid || (mid == num / mid && mid * mid < num)) lo = mid + 1;
            else hi = mid - 1;
        }
        return false;
    }
}
