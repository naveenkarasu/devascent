import java.util.*;
class Solution {
    public long character_replacement(String s, long k) {
        Map<Character, Integer> counts = new HashMap<>();
        int left = 0, maxCount = 0, best = 0;
        for (int right = 0; right < s.length(); right++) {
            char ch = s.charAt(right);
            counts.put(ch, counts.getOrDefault(ch, 0) + 1);
            if (counts.get(ch) > maxCount) maxCount = counts.get(ch);
            while ((right - left + 1) - maxCount > k) {
                char lc = s.charAt(left);
                counts.put(lc, counts.get(lc) - 1);
                left++;
            }
            if (right - left + 1 > best) best = right - left + 1;
        }
        return best;
    }
}
