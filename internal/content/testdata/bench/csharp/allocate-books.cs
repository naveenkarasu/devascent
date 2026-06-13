public class Solution {
    public long allocate_books(long[] pages, long students) {
        int n = pages.Length;
        if (students > n) return -1;
        bool canAllocate(long maxPages) {
            long count = 1, total = 0;
            foreach (long p in pages) {
                if (p > maxPages) return false;
                if (total + p > maxPages) { count++; total = p; }
                else total += p;
            }
            return count <= students;
        }
        long lo = 0, hi = 0;
        foreach (long p in pages) { if (p > lo) lo = p; hi += p; }
        long ans = hi;
        while (lo <= hi) {
            long mid = (lo + hi) / 2;
            if (canAllocate(mid)) { ans = mid; hi = mid - 1; }
            else lo = mid + 1;
        }
        return ans;
    }
}
