class Solution {
    public long count_vowels(String s) {
        long count = 0;
        String vowels = "aeiouAEIOU";
        for (char ch : s.toCharArray()) {
            if (vowels.indexOf(ch) >= 0) count++;
        }
        return count;
    }
}
