import java.util.*;

class Solution {
    public String[] letter_combinations(String digits) {
        if (digits.isEmpty()) return new String[0];
        Map<Character, String> mapping = new HashMap<>();
        mapping.put('2', "abc"); mapping.put('3', "def");
        mapping.put('4', "ghi"); mapping.put('5', "jkl");
        mapping.put('6', "mno"); mapping.put('7', "pqrs");
        mapping.put('8', "tuv"); mapping.put('9', "wxyz");
        List<String> results = new ArrayList<>();
        results.add("");
        for (char d : digits.toCharArray()) {
            String letters = mapping.get(d);
            List<String> next = new ArrayList<>();
            for (String prev : results) {
                for (char ch : letters.toCharArray()) {
                    next.add(prev + ch);
                }
            }
            results = next;
        }
        return results.toArray(new String[0]);
    }
}
