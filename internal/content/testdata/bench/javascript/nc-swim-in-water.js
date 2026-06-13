function swim_in_water(grid) {
  const n = grid.length;
  const visited = new Set();
  // Min-heap: [max_elevation, row, col]
  const heap = [[grid[0][0], 0, 0]];
  visited.add(0);

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
    const [t, i, j] = heapPop(heap);
    if (i === n - 1 && j === n - 1) return t;
    for (const [di, dj] of [[1,0],[-1,0],[0,1],[0,-1]]) {
      const ni = i + di, nj = j + dj;
      if (ni >= 0 && ni < n && nj >= 0 && nj < n && !visited.has(ni * n + nj)) {
        visited.add(ni * n + nj);
        heapPush(heap, [Math.max(t, grid[ni][nj]), ni, nj]);
      }
    }
  }
  return -1;
}
