class Solution {
    public String[] split_odd_even(String s) {
        StringBuilder even = new StringBuilder();
        StringBuilder odd = new StringBuilder();
        for (int i = 0; i < s.length(); i++) {
            if (i % 2 == 0) even.append(s.charAt(i));
            else odd.append(s.charAt(i));
        }
        return new String[]{even.toString(), odd.toString()};
    }
}
