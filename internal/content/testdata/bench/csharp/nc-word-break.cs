using System.Collections.Generic;

public class Solution {
    public bool word_break(string s, string[] word_dict) {
        var wordSet = new HashSet<string>(word_dict);
        int n = s.Length;
        bool[] dp = new bool[n + 1];
        dp[0] = true;
        for (int i = 1; i <= n; i++) {
            for (int j = 0; j < i; j++) {
                if (dp[j] && wordSet.Contains(s.Substring(j, i - j))) {
                    dp[i] = true;
                    break;
                }
            }
        }
        return dp[n];
    }
}
