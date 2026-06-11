import java.util.*;

class Solution {
    public long[][] combination_sum(long[] candidates, long target) {
        Arrays.sort(candidates);
        List<long[]> res = new ArrayList<>();
        backtrack(candidates, 0, target, new ArrayList<>(), res);
        res.sort((a, b) -> {
            int len = Math.min(a.length, b.length);
            for (int i = 0; i < len; i++) {
                if (a[i] != b[i]) return Long.compare(a[i], b[i]);
            }
            return a.length - b.length;
        });
        return res.toArray(new long[0][]);
    }

    private void backtrack(long[] candidates, int start, long remaining, List<Long> current, List<long[]> res) {
        if (remaining == 0) {
            long[] combo = new long[current.size()];
            for (int i = 0; i < current.size(); i++) combo[i] = current.get(i);
            Arrays.sort(combo);
            res.add(combo);
            return;
        }
        for (int i = start; i < candidates.length; i++) {
            if (candidates[i] > remaining) break;
            current.add(candidates[i]);
            backtrack(candidates, i, remaining - candidates[i], current, res);
            current.remove(current.size() - 1);
        }
    }
}
