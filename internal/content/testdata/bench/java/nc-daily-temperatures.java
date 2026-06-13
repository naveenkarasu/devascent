import java.util.*;
class Solution {
    public long[] daily_temperatures(long[] temps) {
        int n = temps.length;
        long[] result = new long[n];
        Deque<int[]> stack = new ArrayDeque<>(); // [temp, index]
        for (int i = 0; i < n; i++) {
            while (!stack.isEmpty() && temps[i] > stack.peek()[0]) {
                int[] top = stack.pop();
                result[top[1]] = i - top[1];
            }
            stack.push(new int[]{(int)temps[i], i});
        }
        return result;
    }
}
