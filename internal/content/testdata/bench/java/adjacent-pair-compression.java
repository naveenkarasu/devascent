import java.util.*;

class Solution {
    public long compressed_length(String s) {
        int n = s.length();
        long count = 0;
        int i = 0;
        while (i < n) {
            if (i + 1 < n && s.charAt(i) == s.charAt(i + 1)) {
                count++;
                i += 2;
            } else {
                count++;
                i++;
            }
        }
        return count;
    }
}
