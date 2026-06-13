function off_bulbs_between_lit(bulbs) {
  const n = bulbs.length;
  let firstOn = -1, lastOn = -1;
  for (let i = 0; i < n; i++) {
    if (bulbs[i] === 1) { firstOn = i; break; }
  }
  for (let i = n - 1; i >= 0; i--) {
    if (bulbs[i] === 1) { lastOn = i; break; }
  }
  if (firstOn === -1) return 0;
  let count = 0;
  for (let j = firstOn; j <= lastOn; j++) {
    if (bulbs[j] === 0) count++;
  }
  return count;
}
