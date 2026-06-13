import java.util.*;
class Solution {
    public long car_fleet(long target, long[] position, long[] speed) {
        int n = position.length;
        int[][] pairs = new int[n][2];
        for (int i = 0; i < n; i++) {
            pairs[i][0] = (int)position[i];
            pairs[i][1] = (int)speed[i];
        }
        Arrays.sort(pairs, (a, b) -> b[0] - a[0]);
        Deque<Double> stack = new ArrayDeque<>();
        for (int[] p : pairs) {
            double time = (double)(target - p[0]) / p[1];
            if (stack.isEmpty() || time > stack.peek()) {
                stack.push(time);
            }
        }
        return stack.size();
    }
}
