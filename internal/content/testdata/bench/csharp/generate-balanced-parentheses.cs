using System.Collections.Generic;

public class Solution {
    public string[] generate_parentheses(long n) {
        var results = new List<string>();
        Backtrack(results, "", 0, 0, (int)n);
        return results.ToArray();
    }

    private void Backtrack(List<string> results, string s, int openCount, int closeCount, int n) {
        if (s.Length == 2 * n) {
            results.Add(s);
            return;
        }
        if (openCount < n) {
            Backtrack(results, s + "(", openCount + 1, closeCount, n);
        }
        if (closeCount < openCount) {
            Backtrack(results, s + ")", openCount, closeCount + 1, n);
        }
    }
}
