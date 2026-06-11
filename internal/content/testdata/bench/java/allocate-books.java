import java.util.*;

class Solution {
    public long allocate_books(long[] pages, long students) {
        int n = pages.length;
        if (students > n) return -1;
        long lo = 0, hi = 0;
        for (long p : pages) {
            lo = Math.max(lo, p);
            hi += p;
        }
        long ans = hi;
        while (lo <= hi) {
            long mid = (lo + hi) / 2;
            if (canAllocate(pages, students, mid)) {
                ans = mid;
                hi = mid - 1;
            } else {
                lo = mid + 1;
            }
        }
        return ans;
    }

    private boolean canAllocate(long[] pages, long students, long maxPages) {
        long count = 1, total = 0;
        for (long p : pages) {
            if (p > maxPages) return false;
            if (total + p > maxPages) {
                count++;
                total = p;
            } else {
                total += p;
            }
        }
        return count <= students;
    }
}
