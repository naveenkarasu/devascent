function max_partition_sum(arr) {
  const total = arr.reduce((a, b) => a + b, 0);
  let current = 0;
  let ans = 0;
  for (let i = 0; i < arr.length - 1; i++) {
    current += arr[i];
    ans = Math.max(ans, Math.max(current, total - current));
  }
  return ans;
}
