import java.util.*;

class Solution {
    public long least_interval(String[] tasks, long n) {
        Map<Character, Long> counts = new HashMap<>();
        for (String t : tasks) {
            char c = t.charAt(0);
            counts.merge(c, 1L, Long::sum);
        }
        long maxCount = Collections.max(counts.values());
        long numMax = counts.values().stream().filter(v -> v == maxCount).count();
        return Math.max(tasks.length, (maxCount - 1) * (n + 1) + numMax);
    }
}
