using System.Collections.Generic;

public class Solution {
    public long[] daily_temperatures(long[] temps) {
        int n = temps.Length;
        long[] result = new long[n];
        var stack = new Stack<int[]>(); // [temp, index]
        for (int i = 0; i < n; i++) {
            while (stack.Count > 0 && temps[i] > stack.Peek()[0]) {
                int[] top = stack.Pop();
                result[top[1]] = i - top[1];
            }
            stack.Push(new int[] { (int)temps[i], i });
        }
        return result;
    }
}
