function can_finish(num_courses, prerequisites) {
  const graph = Array.from({length: num_courses}, () => []);
  for (const [a, b] of prerequisites) {
    graph[a].push(b);
  }
  // 0 = unvisited, 1 = visiting, 2 = done
  const state = new Array(num_courses).fill(0);

  function dfs(c) {
    if (state[c] === 1) return false;
    if (state[c] === 2) return true;
    state[c] = 1;
    for (const nxt of graph[c]) {
      if (!dfs(nxt)) return false;
    }
    state[c] = 2;
    return true;
  }

  for (let c = 0; c < num_courses; c++) {
    if (!dfs(c)) return false;
  }
  return true;
}
