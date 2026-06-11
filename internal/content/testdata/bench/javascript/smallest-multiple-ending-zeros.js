function smallest_trailing_zero_multiple(a, b) {
  const power = Math.pow(10, b);
  function gcd(x, y) {
    while (y !== 0) {
      const t = y;
      y = x % y;
      x = t;
    }
    return x;
  }
  return Math.floor(a / gcd(a, power)) * power;
}
