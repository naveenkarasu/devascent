function selection_sort_k_passes(arr, k) {
  const result = arr.slice();
  const n = result.length;
  const passes = Math.min(k, n - 1);
  for (let i = 0; i < passes; i++) {
    let minIdx = i;
    for (let j = i + 1; j < n; j++) {
      if (result[j] < result[minIdx]) minIdx = j;
    }
    const tmp = result[i];
    result[i] = result[minIdx];
    result[minIdx] = tmp;
  }
  return result;
}
