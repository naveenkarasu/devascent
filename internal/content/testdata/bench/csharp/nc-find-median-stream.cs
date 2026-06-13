using System.Collections.Generic;

public class Solution {
    public object[] median_stream_ops(object[][] operations) {
        // lo = max-heap (negate for PriorityQueue which is min-heap)
        var lo = new PriorityQueue<long, long>(); // max-heap via negative priority
        // hi = min-heap
        var hi = new PriorityQueue<long, long>();
        var out_ = new object[operations.Length];
        for (int i = 0; i < operations.Length; i++) {
            var op = operations[i];
            string name = (string)op[0];
            if (name == "addNum") {
                long v = (long)op[1];
                lo.Enqueue(v, -v); // max-heap: negate priority
                long top = lo.Dequeue();
                hi.Enqueue(top, top);
                if (hi.Count > lo.Count) {
                    long t = hi.Dequeue();
                    lo.Enqueue(t, -t);
                }
                out_[i] = null;
            } else { // findMedian
                double m;
                if (lo.Count > hi.Count) {
                    m = lo.Peek();
                } else {
                    m = (lo.Peek() + hi.Peek()) / 2.0;
                }
                out_[i] = (m == Math.Floor(m)) ? (object)(long)m : (object)m;
            }
        }
        return out_;
    }
}
