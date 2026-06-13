import java.util.*;

class Solution {
    public long[] running_medians(long[] stream) {
        PriorityQueue<Long> lo = new PriorityQueue<>(Collections.reverseOrder()); // max-heap
        PriorityQueue<Long> hi = new PriorityQueue<>(); // min-heap
        long[] result = new long[stream.length];
        for (int idx = 0; idx < stream.length; idx++) {
            long num = stream[idx];
            lo.offer(num);
            hi.offer(lo.poll());
            if (hi.size() > lo.size()) {
                lo.offer(hi.poll());
            }
            if (lo.size() == hi.size()) {
                // both sizes equal: average of two middle values (always whole number in test cases)
                result[idx] = (lo.peek() + hi.peek()) / 2;
            } else {
                result[idx] = lo.peek();
            }
        }
        return result;
    }
}
