import java.util.*;

class Solution {
    public String alien_order(String[] words) {
        Map<Character, Set<Character>> adj = new HashMap<>();
        Map<Character, Integer> indeg = new HashMap<>();
        for (String w : words) {
            for (char c : w.toCharArray()) {
                adj.putIfAbsent(c, new HashSet<>());
                indeg.putIfAbsent(c, 0);
            }
        }
        for (int i = 0; i < words.length - 1; i++) {
            String a = words[i], b = words[i + 1];
            int m = Math.min(a.length(), b.length());
            if (a.length() > b.length() && a.substring(0, m).equals(b.substring(0, m))) {
                return "";
            }
            for (int j = 0; j < m; j++) {
                if (a.charAt(j) != b.charAt(j)) {
                    char from = a.charAt(j), to = b.charAt(j);
                    if (!adj.get(from).contains(to)) {
                        adj.get(from).add(to);
                        indeg.put(to, indeg.get(to) + 1);
                    }
                    break;
                }
            }
        }
        PriorityQueue<Character> heap = new PriorityQueue<>();
        for (char c : indeg.keySet()) {
            if (indeg.get(c) == 0) heap.offer(c);
        }
        StringBuilder res = new StringBuilder();
        while (!heap.isEmpty()) {
            char c = heap.poll();
            res.append(c);
            List<Character> nbrs = new ArrayList<>(adj.get(c));
            Collections.sort(nbrs);
            for (char nxt : nbrs) {
                indeg.put(nxt, indeg.get(nxt) - 1);
                if (indeg.get(nxt) == 0) heap.offer(nxt);
            }
        }
        if (res.length() != indeg.size()) return "";
        return res.toString();
    }
}
