using System.Collections.Generic;

public class Solution {
    public string[] generate_parenthesis(long n) {
        var result = new List<string>();
        Backtrack(result, "", 0, 0, (int)n);
        result.Sort();
        return result.ToArray();
    }
    private void Backtrack(List<string> result, string current, int open, int close, int n) {
        if (current.Length == 2 * n) {
            result.Add(current);
            return;
        }
        if (open < n) Backtrack(result, current + "(", open + 1, close, n);
        if (close < open) Backtrack(result, current + ")", open, close + 1, n);
    }
}
