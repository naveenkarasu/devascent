function network_delay_time(times, n, k) {
  const graph = {};
  for (const [u, v, w] of times) {
    if (!graph[u]) graph[u] = [];
    graph[u].push([w, v]);
  }

  const dist = {};
  // Min-heap: [distance, node]
  const heap = [[0, k]];

  function heapPush(h, item) {
    h.push(item);
    let i = h.length - 1;
    while (i > 0) {
      const parent = Math.floor((i - 1) / 2);
      if (h[parent][0] > h[i][0]) {
        [h[parent], h[i]] = [h[i], h[parent]];
        i = parent;
      } else break;
    }
  }

  function heapPop(h) {
    const top = h[0];
    const last = h.pop();
    if (h.length > 0) {
      h[0] = last;
      let i = 0;
      while (true) {
        const l = 2 * i + 1, r = 2 * i + 2;
        let smallest = i;
        if (l < h.length && h[l][0] < h[smallest][0]) smallest = l;
        if (r < h.length && h[r][0] < h[smallest][0]) smallest = r;
        if (smallest === i) break;
        [h[i], h[smallest]] = [h[smallest], h[i]];
        i = smallest;
      }
    }
    return top;
  }

  while (heap.length > 0) {
    const [d, node] = heapPop(heap);
    if (dist[node] !== undefined) continue;
    dist[node] = d;
    if (graph[node]) {
      for (const [weight, nb] of graph[node]) {
        if (dist[nb] === undefined) {
          heapPush(heap, [d + weight, nb]);
        }
      }
    }
  }

  if (Object.keys(dist).length !== n) return -1;
  return Math.max(...Object.values(dist));
}
