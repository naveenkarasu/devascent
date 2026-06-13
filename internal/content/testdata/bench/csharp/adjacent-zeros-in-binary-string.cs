public class Solution {
    public bool has_adjacent_zeros(string s) {
        for (int i = 1; i < s.Length; i++) {
            if (s[i] == '0' && s[i - 1] == '0') return true;
        }
        return false;
    }
}
