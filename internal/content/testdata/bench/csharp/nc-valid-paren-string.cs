public class Solution {
    public bool check_valid_string(string s) {
        long lo = 0, hi = 0;
        foreach (char c in s) {
            if (c == '(') { lo++; hi++; }
            else if (c == ')') { lo--; hi--; }
            else { lo--; hi++; }
            if (hi < 0) return false;
            if (lo < 0) lo = 0;
        }
        return lo == 0;
    }
}
