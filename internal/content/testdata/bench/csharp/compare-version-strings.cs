public class Solution {
    public long compare_versions(string v1, string v2) {
        string[] parts1 = v1.Split('.');
        string[] parts2 = v2.Split('.');
        int length = Math.Max(parts1.Length, parts2.Length);
        for (int i = 0; i < length; i++) {
            long a = (i < parts1.Length) ? long.Parse(parts1[i]) : 0;
            long b = (i < parts2.Length) ? long.Parse(parts2[i]) : 0;
            if (a > b) return 1;
            if (a < b) return -1;
        }
        return 0;
    }
}
