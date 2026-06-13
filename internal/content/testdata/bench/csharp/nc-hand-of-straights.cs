using System.Collections.Generic;

public class Solution {
    public bool is_n_straight_hand(long[] hand, long group_size) {
        if (hand.Length % group_size != 0) return false;
        var count = new SortedDictionary<long, long>();
        foreach (long card in hand) {
            if (!count.ContainsKey(card)) count[card] = 0;
            count[card]++;
        }
        foreach (long card in new List<long>(count.Keys)) {
            long need = count.ContainsKey(card) ? count[card] : 0;
            if (need > 0) {
                for (long i = 0; i < group_size; i++) {
                    long c = count.ContainsKey(card + i) ? count[card + i] : 0;
                    if (c < need) return false;
                    count[card + i] = c - need;
                }
            }
        }
        return true;
    }
}
