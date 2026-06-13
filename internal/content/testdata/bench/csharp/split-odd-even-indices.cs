public class Solution {
    public string[] split_odd_even(string s) {
        var even = new System.Text.StringBuilder();
        var odd = new System.Text.StringBuilder();
        for (int i = 0; i < s.Length; i++) {
            if (i % 2 == 0) even.Append(s[i]);
            else odd.Append(s[i]);
        }
        return new string[] { even.ToString(), odd.ToString() };
    }
}
