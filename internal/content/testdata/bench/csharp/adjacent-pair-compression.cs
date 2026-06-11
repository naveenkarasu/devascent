public class Solution {
    public long compressed_length(string s) {
        int n = s.Length;
        long count = 0;
        int i = 0;
        while (i < n) {
            if (i + 1 < n && s[i] == s[i + 1]) {
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
