public class Solution {
    public long min_changes_no_triple(string s) {
        var arr = s.ToCharArray();
        long ans = 0;
        int same = 1;
        for (int i = 1; i < arr.Length; i++) {
            if (arr[i] == arr[i - 1]) {
                same++;
            } else {
                same = 1;
            }
            if (same == 3) {
                ans++;
                arr[i] = '@';
                same = 1;
            }
        }
        return ans;
    }
}
