function triple_bucket_sum(counts) {
  const MOD = 1000000007n;
  const n = counts.length;
  let ans = 0n;
  for (let i = 0; i < n - 2; i++) {
    for (let j = i + 1; j < n - 1; j++) {
      for (let k = j + 1; k < n; k++) {
        ans = (ans + BigInt(counts[i]) * BigInt(counts[j]) * BigInt(counts[k])) % MOD;
      }
    }
  }
  return Number(ans);
}
