function is_fibonacci(n) {
  if (n <= 0) return false;
  let a = 1, b = 1;
  while (a < n) {
    [a, b] = [b, a + b];
  }
  return a === n;
}
