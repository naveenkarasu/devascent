using System.Collections.Generic;

public class Solution {
    public string[] letter_combinations(string digits) {
        if (digits.Length == 0) return new string[0];
        var mapping = new Dictionary<char, string> {
            {'2', "abc"}, {'3', "def"}, {'4', "ghi"}, {'5', "jkl"},
            {'6', "mno"}, {'7', "pqrs"}, {'8', "tuv"}, {'9', "wxyz"}
        };
        var results = new List<string> { "" };
        foreach (char d in digits) {
            string letters = mapping[d];
            var next = new List<string>();
            foreach (string prev in results) {
                foreach (char ch in letters) {
                    next.Add(prev + ch);
                }
            }
            results = next;
        }
        return results.ToArray();
    }
}
