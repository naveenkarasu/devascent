import java.util.*;

class Solution {
    public long steps_to_overtake(long a, long b) {
        long steps = 0;
        while (a <= b) {
            a *= 3;
            b *= 2;
            steps++;
        }
        return steps;
    }
}
