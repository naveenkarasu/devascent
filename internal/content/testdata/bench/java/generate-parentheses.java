import java.util.*;

class Solution {
    public String[] generate_parentheses(long n) {
        List<String> result = new ArrayList<>();
        backtrack(result, "", 0, 0, (int) n);
        Collections.sort(result);
        return result.toArray(new String[0]);
    }

    private void backtrack(List<String> result, String s, int open, int close, int n) {
        if (s.length() == 2 * n) {
            result.add(s);
            return;
        }
        if (open < n) backtrack(result, s + "(", open + 1, close, n);
        if (close < open) backtrack(result, s + ")", open, close + 1, n);
    }
}
