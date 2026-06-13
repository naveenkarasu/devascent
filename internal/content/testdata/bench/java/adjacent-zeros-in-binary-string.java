class Solution {
    public boolean has_adjacent_zeros(String s) {
        for (int i = 1; i < s.length(); i++) {
            if (s.charAt(i) == '0' && s.charAt(i - 1) == '0') return true;
        }
        return false;
    }
}
