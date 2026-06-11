using System.Collections.Generic;

public class Solution {
    public double[] running_medians(long[] stream) {
        // lo is max-heap (negate priority), hi is min-heap
        var lo = new PriorityQueue<long, long>(); // max-heap via negative priority
        var hi = new PriorityQueue<long, long>(); // min-heap
        var result = new double[stream.Length];
        for (int idx = 0; idx < stream.Length; idx++) {
            long num = stream[idx];
            lo.Enqueue(num, -num);
            long top = lo.Dequeue();
            hi.Enqueue(top, top);
            if (hi.Count > lo.Count) {
                long val = hi.Dequeue();
                lo.Enqueue(val, -val);
            }
            if (lo.Count == hi.Count) {
                result[idx] = (lo.Peek() + hi.Peek()) / 2.0;
            } else {
                result[idx] = (double)lo.Peek();
            }
        }
        return result;
    }
}
