function perfect_numbers_up_to(n) {
  const result = [];
  for (let i = 2; i <= n; i++) {
    let s = 1;
    let r = Math.floor(Math.sqrt(i));
    while (r * r > i) r--;
    while ((r + 1) * (r + 1) <= i) r++;
    for (let d = 2; d <= r; d++) {
      if (i % d === 0) {
        s += d + Math.floor(i / d);
        if (d * d === i) s -= d;
      }
    }
    if (s === i) result.push(i);
  }
  return result;
}
