import java.util.*;

class Solution {
    public boolean is_anagram(String s, String t) {
        if (s.length() != t.length()) return false;
        Map<Character, Integer> counts = new HashMap<>();
        for (char c : s.toCharArray()) counts.merge(c, 1, Integer::sum);
        for (char c : t.toCharArray()) {
            if (counts.getOrDefault(c, 0) == 0) return false;
            counts.merge(c, -1, Integer::sum);
        }
        return true;
    }
}
