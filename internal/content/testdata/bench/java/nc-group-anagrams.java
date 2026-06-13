import java.util.*;

class Solution {
    public String[][] group_anagrams(String[] strs) {
        Map<String, List<String>> groups = new HashMap<>();
        for (String s : strs) {
            char[] arr = s.toCharArray();
            Arrays.sort(arr);
            String key = new String(arr);
            groups.computeIfAbsent(key, k -> new ArrayList<>()).add(s);
        }
        List<List<String>> res = new ArrayList<>();
        for (List<String> g : groups.values()) {
            Collections.sort(g);
            res.add(g);
        }
        res.sort((a, b) -> {
            for (int i = 0; i < Math.min(a.size(), b.size()); i++) {
                int cmp = a.get(i).compareTo(b.get(i));
                if (cmp != 0) return cmp;
            }
            return a.size() - b.size();
        });
        String[][] out = new String[res.size()][];
        for (int i = 0; i < res.size(); i++) {
            out[i] = res.get(i).toArray(new String[0]);
        }
        return out;
    }
}
