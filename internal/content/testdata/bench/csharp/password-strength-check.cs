public class Solution {
    public bool is_strong_password(string s) {
        if (s.Length < 8) return false;
        bool hasUpper = false, hasLower = false, hasDigit = false, hasSpecial = false;
        foreach (char c in s) {
            if (char.IsUpper(c)) hasUpper = true;
            if (char.IsLower(c)) hasLower = true;
            if (char.IsDigit(c)) hasDigit = true;
            if (!char.IsLetterOrDigit(c)) hasSpecial = true;
        }
        return hasUpper && hasLower && hasDigit && hasSpecial;
    }
}
