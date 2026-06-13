import java.util.*;

class Solution {
    public String[] letter_combinations(String digits) {
        if (digits.isEmpty()) return new String[0];
        String[] mapping = {"", "", "abc", "def", "ghi", "jkl", "mno", "pqrs", "tuv", "wxyz"};
        List<String> res = new ArrayList<>();
        backtrack(digits, 0, "", mapping, res);
        Collections.sort(res);
        return res.toArray(new String[0]);
    }

    private void backtrack(String digits, int index, String current, String[] mapping, List<String> res) {
        if (index == digits.length()) {
            res.add(current);
            return;
        }
        for (char ch : mapping[digits.charAt(index) - '0'].toCharArray()) {
            backtrack(digits, index + 1, current + ch, mapping, res);
        }
    }
}
