using System.Collections.Generic;

public class Solution {
    public object[] min_stack_ops(object[][] operations) {
        var stack = new Stack<long>();
        var mins = new Stack<long>();
        var out_ = new object[operations.Length];
        for (int i = 0; i < operations.Length; i++) {
            var op = operations[i];
            string name = (string)op[0];
            switch (name) {
                case "push": {
                    long v = (long)op[1];
                    stack.Push(v);
                    mins.Push(mins.Count == 0 ? v : Math.Min(v, mins.Peek()));
                    out_[i] = null;
                    break;
                }
                case "pop":
                    stack.Pop();
                    mins.Pop();
                    out_[i] = null;
                    break;
                case "top":
                    out_[i] = stack.Peek();
                    break;
                default: // getMin
                    out_[i] = mins.Peek();
                    break;
            }
        }
        return out_;
    }
}
