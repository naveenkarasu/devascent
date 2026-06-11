import java.util.*;

class Solution {
    public String[] generate_parentheses(long n) {
        List<String> results = new ArrayList<>();
        backtrack(results, "", 0, 0, (int) n);
        return results.toArray(new String[0]);
    }

    private void backtrack(List<String> results, String s, int openCount, int closeCount, int n) {
        if (s.length() == 2 * n) {
            results.add(s);
            return;
        }
        if (openCount < n) backtrack(results, s + "(", openCount + 1, closeCount, n);
        if (closeCount < openCount) backtrack(results, s + ")", openCount, closeCount + 1, n);
    }
}
