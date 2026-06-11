using System;
using System.Collections.Generic;

public class Solution {
    public long car_fleet(long target, long[] position, long[] speed) {
        int n = position.Length;
        var pairs = new (long pos, long spd)[n];
        for (int i = 0; i < n; i++) pairs[i] = (position[i], speed[i]);
        Array.Sort(pairs, (a, b) => b.pos.CompareTo(a.pos));
        var stack = new Stack<double>();
        foreach (var (pos, spd) in pairs) {
            double time = (double)(target - pos) / spd;
            if (stack.Count == 0 || time > stack.Peek()) {
                stack.Push(time);
            }
        }
        return stack.Count;
    }
}
