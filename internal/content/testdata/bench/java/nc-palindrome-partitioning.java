import java.util.*;

class Solution {
    public String[][] partition(String s) {
        List<String[]> res = new ArrayList<>();
        backtrack(s, 0, new ArrayList<>(), res);
        res.sort((a, b) -> {
            int len = Math.min(a.length, b.length);
            for (int i = 0; i < len; i++) {
                int cmp = a[i].compareTo(b[i]);
                if (cmp != 0) return cmp;
            }
            return Integer.compare(a.length, b.length);
        });
        return res.toArray(new String[0][]);
    }

    private void backtrack(String s, int start, List<String> current, List<String[]> res) {
        if (start == s.length()) {
            res.add(current.toArray(new String[0]));
            return;
        }
        for (int end = start + 1; end <= s.length(); end++) {
            String substr = s.substring(start, end);
            if (isPalindrome(substr)) {
                current.add(substr);
                backtrack(s, end, current, res);
                current.remove(current.size() - 1);
            }
        }
    }

    private boolean isPalindrome(String t) {
        int l = 0, r = t.length() - 1;
        while (l < r) {
            if (t.charAt(l) != t.charAt(r)) return false;
            l++; r--;
        }
        return true;
    }
}
