import java.util.*;
class Solution {
    public boolean check_inclusion(String s1, String s2) {
        if (s1.length() > s2.length()) return false;
        Map<Character, Integer> need = new HashMap<>();
        for (char c : s1.toCharArray()) need.put(c, need.getOrDefault(c, 0) + 1);
        Map<Character, Integer> window = new HashMap<>();
        int k = s1.length();
        for (int i = 0; i < s2.length(); i++) {
            char ch = s2.charAt(i);
            window.put(ch, window.getOrDefault(ch, 0) + 1);
            if (i >= k) {
                char lc = s2.charAt(i - k);
                window.put(lc, window.get(lc) - 1);
                if (window.get(lc) == 0) window.remove(lc);
            }
            if (window.equals(need)) return true;
        }
        return false;
    }
}
