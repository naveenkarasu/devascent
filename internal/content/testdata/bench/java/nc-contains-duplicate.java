import java.util.*;

class Solution {
    public boolean contains_duplicate(long[] nums) {
        Set<Long> seen = new HashSet<>();
        for (long n : nums) {
            if (!seen.add(n)) return true;
        }
        return false;
    }
}
