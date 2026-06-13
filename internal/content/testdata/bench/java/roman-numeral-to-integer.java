import java.util.*;

class Solution {
    public long roman_to_int(String s) {
        Map<String, Long> doubles = new HashMap<>();
        doubles.put("CM", 900L); doubles.put("CD", 400L);
        doubles.put("XC", 90L); doubles.put("XL", 40L);
        doubles.put("IX", 9L); doubles.put("IV", 4L);
        Map<Character, Long> singles = new HashMap<>();
        singles.put('M', 1000L); singles.put('D', 500L);
        singles.put('C', 100L); singles.put('L', 50L);
        singles.put('X', 10L); singles.put('V', 5L);
        singles.put('I', 1L);
        long total = 0;
        int i = 0;
        while (i < s.length()) {
            if (i < s.length() - 1 && doubles.containsKey(s.substring(i, i + 2))) {
                total += doubles.get(s.substring(i, i + 2));
                i += 2;
            } else {
                total += singles.get(s.charAt(i));
                i++;
            }
        }
        return total;
    }
}
