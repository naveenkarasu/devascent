using System.Collections.Generic;

public class Solution {
    public long eval_rpn(string[] tokens) {
        var stack = new Stack<long>();
        foreach (string t in tokens) {
            if (t == "+" || t == "-" || t == "*" || t == "/") {
                long b = stack.Pop(), a = stack.Pop();
                if (t == "+") stack.Push(a + b);
                else if (t == "-") stack.Push(a - b);
                else if (t == "*") stack.Push(a * b);
                else stack.Push((long)(a / b));
            } else {
                stack.Push(long.Parse(t));
            }
        }
        return stack.Pop();
    }
}
