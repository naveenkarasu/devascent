import java.util.*;

class Solution {
    public Object[] median_stream_ops(Object[][] operations) {
        PriorityQueue<Long> lo = new PriorityQueue<>(Collections.reverseOrder()); // max-heap
        PriorityQueue<Long> hi = new PriorityQueue<>();                            // min-heap
        Object[] out = new Object[operations.length];
        for (int i = 0; i < operations.length; i++) {
            Object[] op = operations[i];
            String name = (String) op[0];
            if (name.equals("addNum")) {
                long v = ((Number) op[1]).longValue();
                lo.add(v);
                hi.add(lo.poll());
                if (hi.size() > lo.size()) lo.add(hi.poll());
                out[i] = null;
            } else { // findMedian
                double m;
                if (lo.size() > hi.size()) {
                    m = lo.peek();
                } else {
                    m = (lo.peek() + hi.peek()) / 2.0;
                }
                // value-based boxing: integral -> Long, fractional -> Double
                out[i] = (m == Math.rint(m)) ? (Object) (long) m : (Object) m;
            }
        }
        return out;
    }
}
