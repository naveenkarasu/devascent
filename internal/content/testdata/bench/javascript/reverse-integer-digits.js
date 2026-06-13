function reverse_integer(n) {
  let reversed_n = 0;
  while (n > 0) {
    reversed_n = reversed_n * 10 + (n % 10);
    n = Math.floor(n / 10);
  }
  return reversed_n;
}
