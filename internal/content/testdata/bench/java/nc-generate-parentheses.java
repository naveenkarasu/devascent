import java.util.*;
class Solution {
    public String[] generate_parenthesis(long n) {
        List<String> result = new ArrayList<>();
        backtrack(result, "", 0, 0, (int)n);
        Collections.sort(result);
        return result.toArray(new String[0]);
    }
    private void backtrack(List<String> result, String current, int open, int close, int n) {
        if (current.length() == 2 * n) {
            result.add(current);
            return;
        }
        if (open < n) backtrack(result, current + "(", open + 1, close, n);
        if (close < open) backtrack(result, current + ")", open, close + 1, n);
    }
}
