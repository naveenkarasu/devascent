using System.Collections.Generic;

public class Solution {
    public string[] generate_parentheses(long n) {
        var result = new List<string>();
        void backtrack(string s, int open, int close) {
            if (s.Length == 2 * (int)n) { result.Add(s); return; }
            if (open < (int)n) backtrack(s + "(", open + 1, close);
            if (close < open) backtrack(s + ")", open, close + 1);
        }
        backtrack("", 0, 0);
        result.Sort();
        return result.ToArray();
    }
}
