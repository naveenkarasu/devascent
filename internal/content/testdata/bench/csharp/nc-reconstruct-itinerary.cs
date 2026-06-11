public class Solution {
    public string[] find_itinerary(string[][] tickets) {
        var graph = new Dictionary<string, List<string>>();
        // Sort tickets in reverse so that RemoveLast() yields lexicographically smallest
        var sortedTickets = tickets.OrderBy(t => t[0]).ThenBy(t => t[1]).Reverse().ToArray();
        foreach (var ticket in sortedTickets) {
            if (!graph.ContainsKey(ticket[0])) graph[ticket[0]] = new List<string>();
            graph[ticket[0]].Add(ticket[1]);
        }
        var result = new List<string>();
        Dfs("JFK", graph, result);
        result.Reverse();
        return result.ToArray();
    }

    private void Dfs(string airport, Dictionary<string, List<string>> graph, List<string> result) {
        var neighbors = graph.ContainsKey(airport) ? graph[airport] : new List<string>();
        while (neighbors.Count > 0) {
            string next = neighbors[neighbors.Count - 1];
            neighbors.RemoveAt(neighbors.Count - 1);
            Dfs(next, graph, result);
        }
        result.Add(airport);
    }
}
