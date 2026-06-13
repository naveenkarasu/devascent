function free_vertices(n, m) {
  let k = 0;
  while (Math.floor(k * (k - 1) / 2) < m) {
    k++;
  }
  return n - k;
}
