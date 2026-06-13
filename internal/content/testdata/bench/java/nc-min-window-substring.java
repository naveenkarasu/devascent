import java.util.*;
class Solution {
    public String min_window(String s, String t) {
        if (t.isEmpty() || s.isEmpty()) return "";
        Map<Character, Integer> need = new HashMap<>();
        for (char c : t.toCharArray()) need.put(c, need.getOrDefault(c, 0) + 1);
        Map<Character, Integer> have = new HashMap<>();
        int formed = 0, required = need.size();
        int left = 0, bestLen = Integer.MAX_VALUE, bestStart = 0;
        for (int right = 0; right < s.length(); right++) {
            char ch = s.charAt(right);
            have.put(ch, have.getOrDefault(ch, 0) + 1);
            if (need.containsKey(ch) && have.get(ch).equals(need.get(ch))) formed++;
            while (formed == required) {
                if (right - left + 1 < bestLen) { bestLen = right - left + 1; bestStart = left; }
                char lc = s.charAt(left);
                have.put(lc, have.get(lc) - 1);
                if (need.containsKey(lc) && have.get(lc) < need.get(lc)) formed--;
                left++;
            }
        }
        return bestLen == Integer.MAX_VALUE ? "" : s.substring(bestStart, bestStart + bestLen);
    }
}
