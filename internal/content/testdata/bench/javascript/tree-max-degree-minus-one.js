function trap_node(n, edges) {
  const degree = new Array(n).fill(0);
  for (const [u, v] of edges) {
    degree[u]++;
    degree[v]++;
  }
  return Math.max(...degree) - 1;
}
