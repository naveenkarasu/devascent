import java.util.*;

class Solution {
    public boolean is_n_straight_hand(long[] hand, long group_size) {
        if (hand.length % group_size != 0) return false;
        TreeMap<Long, Long> count = new TreeMap<>();
        for (long card : hand) count.merge(card, 1L, Long::sum);
        for (long card : new ArrayList<>(count.keySet())) {
            long need = count.getOrDefault(card, 0L);
            if (need > 0) {
                for (long i = 0; i < group_size; i++) {
                    long c = count.getOrDefault(card + i, 0L);
                    if (c < need) return false;
                    count.put(card + i, c - need);
                }
            }
        }
        return true;
    }
}
