function edit_distance(str1, str2) {
  const m = str1.length, n = str2.length;
  let cur = Array.from({length: n + 1}, (_, i) => i);
  for (let i = 1; i <= m; i++) {
    let pre = cur[0];
    cur[0] = i;
    for (let j = 1; j <= n; j++) {
      const temp = cur[j];
      if (str1[i - 1] === str2[j - 1]) {
        cur[j] = pre;
      } else {
        cur[j] = Math.min(pre, cur[j - 1], cur[j]) + 1;
      }
      pre = temp;
    }
  }
  return cur[n];
}
