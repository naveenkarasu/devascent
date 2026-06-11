import java.util.*;

class Solution {
    public boolean merge_triplets(long[][] triplets, long[] target) {
        long[] result = {0, 0, 0};
        for (long[] t : triplets) {
            if (t[0] <= target[0] && t[1] <= target[1] && t[2] <= target[2]) {
                result[0] = Math.max(result[0], t[0]);
                result[1] = Math.max(result[1], t[1]);
                result[2] = Math.max(result[2], t[2]);
            }
        }
        return result[0] == target[0] && result[1] == target[1] && result[2] == target[2];
    }
}
