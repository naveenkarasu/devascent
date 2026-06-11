function min_cost_connect_points(points) {
  const n = points.length;
  const visited = new Set();
  // Min-heap: [cost, point_index]
  // Using a simple array-based min-heap
  const heap = [[0, 0]];
  let total = 0;

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

  while (visited.size < n) {
    const [cost, i] = heapPop(heap);
    if (visited.has(i)) continue;
    visited.add(i);
    total += cost;
    const [xi, yi] = points[i];
    for (let j = 0; j < n; j++) {
      if (!visited.has(j)) {
        const [xj, yj] = points[j];
        heapPush(heap, [Math.abs(xi - xj) + Math.abs(yi - yj), j]);
      }
    }
  }
  return total;
}
