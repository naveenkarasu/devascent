using System.Collections.Generic;

public class Solution {
    public bool is_valid(string s) {
        var stack = new Stack<char>();
        foreach (char ch in s) {
            if (ch == ')' || ch == ']' || ch == '}') {
                if (stack.Count == 0) return false;
                char top = stack.Pop();
                if (ch == ')' && top != '(') return false;
                if (ch == ']' && top != '[') return false;
                if (ch == '}' && top != '{') return false;
            } else {
                stack.Push(ch);
            }
        }
        return stack.Count == 0;
    }
}
