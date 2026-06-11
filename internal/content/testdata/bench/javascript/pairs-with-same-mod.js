function count_pairs_same_mod(nums, divisor) {
  const freq = {};
  for (const x of nums) {
    const r = ((x % divisor) + divisor) % divisor;
    freq[r] = (freq[r] || 0) + 1;
  }
  let pairs = 0;
  for (const cnt of Object.values(freq)) {
    pairs += Math.floor(cnt * (cnt - 1) / 2);
  }
  return pairs;
}
