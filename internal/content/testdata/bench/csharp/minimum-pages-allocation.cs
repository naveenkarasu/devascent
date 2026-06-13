public class Solution {
    public long min_max_pages(long[] pages, long num_students) {
        if (num_students > pages.Length) return -1;
        long lo = 0, hi = 0;
        foreach (long p in pages) {
            if (p > lo) lo = p;
            hi += p;
        }
        long result = hi;
        while (lo <= hi) {
            long mid = (lo + hi) / 2;
            if (CanAllocate(pages, num_students, mid)) {
                result = mid;
                hi = mid - 1;
            } else {
                lo = mid + 1;
            }
        }
        return result;
    }

    private bool CanAllocate(long[] pages, long numStudents, long maxPages) {
        long students = 1, current = 0;
        foreach (long p in pages) {
            if (p > maxPages) return false;
            if (current + p > maxPages) {
                students++;
                current = p;
            } else {
                current += p;
            }
        }
        return students <= numStudents;
    }
}
