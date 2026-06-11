import java.util.*;

class Solution {
    public Object[] min_stack_ops(Object[][] operations) {
        Deque<Long> stack = new ArrayDeque<>();
        Deque<Long> mins = new ArrayDeque<>();
        Object[] out = new Object[operations.length];
        for (int i = 0; i < operations.length; i++) {
            Object[] op = operations[i];
            String name = (String) op[0];
            switch (name) {
                case "push": {
                    long v = ((Number) op[1]).longValue();
                    stack.push(v);
                    mins.push(mins.isEmpty() ? v : Math.min(v, mins.peek()));
                    out[i] = null;
                    break;
                }
                case "pop":
                    stack.pop();
                    mins.pop();
                    out[i] = null;
                    break;
                case "top":
                    out[i] = stack.peek();
                    break;
                default: // getMin
                    out[i] = mins.peek();
                    break;
            }
        }
        return out;
    }
}
