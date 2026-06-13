class Solution {
    public long count_common_characters(String[] strings) {
        if (strings.length == 0) return 0;
        java.util.Set<Character> common = new java.util.HashSet<>();
        for (char c : strings[0].toCharArray()) common.add(c);
        for (int i = 1; i < strings.length; i++) {
            java.util.Set<Character> cur = new java.util.HashSet<>();
            for (char c : strings[i].toCharArray()) cur.add(c);
            common.retainAll(cur);
        }
        return common.size();
    }
}
