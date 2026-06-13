function count_modeq_pairs(n, m) {
  let count = 0;
  for (let a = 1; a <= n; a++) {
    for (let b = a + 1; b <= n; b++) {
      if ((m % a) % b === (m % b) % a) count++;
    }
  }
  return count;
}
