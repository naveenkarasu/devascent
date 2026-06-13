function product_of_counts(s) {
  const MOD = 1000000007n;
  const freq = {};
  for (const ch of s) {
    freq[ch] = (freq[ch] || 0) + 1;
  }
  let result = 1n;
  for (const cnt of Object.values(freq)) {
    result = (result * BigInt(cnt)) % MOD;
  }
  return Number(result);
}
