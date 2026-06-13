import java.util.*;
class Solution {
    public long eval_rpn(String[] tokens) {
        Deque<Long> stack = new ArrayDeque<>();
        Set<String> ops = new HashSet<>(Arrays.asList("+", "-", "*", "/"));
        for (String t : tokens) {
            if (ops.contains(t)) {
                long b = stack.pop(), a = stack.pop();
                if (t.equals("+")) stack.push(a + b);
                else if (t.equals("-")) stack.push(a - b);
                else if (t.equals("*")) stack.push(a * b);
                else stack.push((long)(a / b));
            } else {
                stack.push(Long.parseLong(t));
            }
        }
        return stack.pop();
    }
}
