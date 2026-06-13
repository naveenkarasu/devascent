function sum_squared_digits(n) {
  let total = 0;
  while (n > 0) {
    const d = n % 10;
    total += d * d;
    n = Math.floor(n / 10);
  }
  return total;
}
