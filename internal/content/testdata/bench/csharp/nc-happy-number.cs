using System.Collections.Generic;

public class Solution {
    public bool is_happy(long n) {
        var seen = new HashSet<long>();
        while (n != 1) {
            if (seen.Contains(n)) return false;
            seen.Add(n);
            long total = 0;
            while (n > 0) {
                long digit = n % 10;
                total += digit * digit;
                n /= 10;
            }
            n = total;
        }
        return true;
    }
}
