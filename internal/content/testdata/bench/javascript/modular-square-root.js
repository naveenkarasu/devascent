function modular_sqrt(n, m) {
  for (let x = 0; x < m; x++) {
    if ((x * x) % m === n) return x;
  }
  return -1;
}
