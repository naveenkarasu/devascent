import java.util.*;

class Solution {
    public long[] rectangle_perimeter_area(long a, long b) {
        if (a <= 0 || b <= 0) return new long[]{0, 0};
        return new long[]{(a + b) * 2, a * b};
    }
}
