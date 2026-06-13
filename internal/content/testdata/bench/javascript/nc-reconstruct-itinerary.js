function find_itinerary(tickets) {
  const graph = {};
  // Sort descending so pop() yields smallest first
  const sorted = tickets.slice().sort((a, b) => {
    if (a[0] > b[0]) return -1;
    if (a[0] < b[0]) return 1;
    if (a[1] > b[1]) return -1;
    if (a[1] < b[1]) return 1;
    return 0;
  });
  for (const [src, dst] of sorted) {
    if (!graph[src]) graph[src] = [];
    graph[src].push(dst);
  }

  const result = [];
  function dfs(airport) {
    while (graph[airport] && graph[airport].length > 0) {
      dfs(graph[airport].pop());
    }
    result.push(airport);
  }

  dfs('JFK');
  return result.reverse();
}
