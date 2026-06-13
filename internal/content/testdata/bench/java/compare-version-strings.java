class Solution {
    public long compare_versions(String v1, String v2) {
        String[] parts1 = v1.split("\\.");
        String[] parts2 = v2.split("\\.");
        int length = Math.max(parts1.length, parts2.length);
        for (int i = 0; i < length; i++) {
            long a = (i < parts1.length) ? Long.parseLong(parts1[i]) : 0;
            long b = (i < parts2.length) ? Long.parseLong(parts2[i]) : 0;
            if (a > b) return 1;
            else if (a < b) return -1;
        }
        return 0;
    }
}
