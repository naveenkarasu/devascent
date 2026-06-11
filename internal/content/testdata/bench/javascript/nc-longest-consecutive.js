function longest_consecutive(nums) {
  const s = new Set(nums);
  let best = 0;
  for (const n of s) {
    if (!s.has(n - 1)) {
      let length = 1;
      while (s.has(n + length)) length++;
      best = Math.max(best, length);
    }
  }
  return best;
}
