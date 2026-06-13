import java.util.*;

class Solution {
    public String[] find_itinerary(String[][] tickets) {
        Map<String, List<String>> graph = new HashMap<>();
        // Sort tickets in reverse so that pop() gives lexicographically smallest
        Arrays.sort(tickets, (a, b) -> {
            int c = a[0].compareTo(b[0]);
            if (c != 0) return -c;
            return -a[1].compareTo(b[1]);
        });
        for (String[] ticket : tickets) {
            graph.computeIfAbsent(ticket[0], k -> new ArrayList<>()).add(ticket[1]);
        }
        List<String> result = new ArrayList<>();
        dfs("JFK", graph, result);
        Collections.reverse(result);
        return result.toArray(new String[0]);
    }

    private void dfs(String airport, Map<String, List<String>> graph, List<String> result) {
        List<String> neighbors = graph.getOrDefault(airport, new ArrayList<>());
        while (!neighbors.isEmpty()) {
            String next = neighbors.remove(neighbors.size() - 1);
            dfs(next, graph, result);
        }
        result.add(airport);
    }
}
