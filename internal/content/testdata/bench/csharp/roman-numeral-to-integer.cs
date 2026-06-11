using System.Collections.Generic;

public class Solution {
    public long roman_to_int(string s) {
        var doubles = new Dictionary<string, long> {
            {"CM", 900}, {"CD", 400}, {"XC", 90}, {"XL", 40}, {"IX", 9}, {"IV", 4}
        };
        var singles = new Dictionary<char, long> {
            {'M', 1000}, {'D', 500}, {'C', 100}, {'L', 50}, {'X', 10}, {'V', 5}, {'I', 1}
        };
        long total = 0;
        int i = 0;
        while (i < s.Length) {
            if (i < s.Length - 1 && doubles.TryGetValue(s.Substring(i, 2), out long dval)) {
                total += dval;
                i += 2;
            } else {
                total += singles[s[i]];
                i++;
            }
        }
        return total;
    }
}
