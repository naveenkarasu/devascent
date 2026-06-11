class Solution {
    public String multiply_strings(String num1, String num2) {
        int m = num1.length(), n = num2.length();
        int[] pos = new int[m + n];
        for (int i = m - 1; i >= 0; i--) {
            for (int j = n - 1; j >= 0; j--) {
                int mul = (num1.charAt(i) - '0') * (num2.charAt(j) - '0');
                int p1 = i + j, p2 = i + j + 1;
                int total = mul + pos[p2];
                pos[p2] = total % 10;
                pos[p1] += total / 10;
            }
        }
        StringBuilder sb = new StringBuilder();
        for (int d : pos) {
            if (sb.length() == 0 && d == 0) continue;
            sb.append(d);
        }
        return sb.length() == 0 ? "0" : sb.toString();
    }
}
