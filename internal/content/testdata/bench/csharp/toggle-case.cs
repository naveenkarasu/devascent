public class Solution {
    public string toggle_case(string s) {
        var result = new System.Text.StringBuilder();
        foreach (char ch in s) {
            if (char.IsUpper(ch)) result.Append(char.ToLower(ch));
            else if (char.IsLower(ch)) result.Append(char.ToUpper(ch));
            else result.Append(ch);
        }
        return result.ToString();
    }
}
