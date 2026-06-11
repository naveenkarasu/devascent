import java.util.*;
class Solution {
    public long length_of_longest_substring(String s) {
        Map<Character, Integer> lastSeen = new HashMap<>();
        int left = 0, best = 0;
        for (int right = 0; right < s.length(); right++) {
            char ch = s.charAt(right);
            if (lastSeen.containsKey(ch) && lastSeen.get(ch) >= left)
                left = lastSeen.get(ch) + 1;
            lastSeen.put(ch, right);
            if (right - left + 1 > best) best = right - left + 1;
        }
        return best;
    }
}
