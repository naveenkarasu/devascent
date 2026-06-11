public class Solution {
    public string[] letter_combinations(string digits) {
        if (digits.Length == 0) return new string[0];
        string[] mapping = {"", "", "abc", "def", "ghi", "jkl", "mno", "pqrs", "tuv", "wxyz"};
        var res = new List<string>();
        Backtrack(digits, 0, "", mapping, res);
        res.Sort();
        return res.ToArray();
    }

    private void Backtrack(string digits, int index, string current, string[] mapping, List<string> res) {
        if (index == digits.Length) {
            res.Add(current);
            return;
        }
        foreach (char ch in mapping[digits[index] - '0']) {
            Backtrack(digits, index + 1, current + ch, mapping, res);
        }
    }
}
