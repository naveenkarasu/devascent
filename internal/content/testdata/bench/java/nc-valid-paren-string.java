import java.util.*;

class Solution {
    public boolean check_valid_string(String s) {
        long lo = 0, hi = 0;
        for (char c : s.toCharArray()) {
            if (c == '(') {
                lo++; hi++;
            } else if (c == ')') {
                lo--; hi--;
            } else {
                lo--; hi++;
            }
            if (hi < 0) return false;
            if (lo < 0) lo = 0;
        }
        return lo == 0;
    }
}
