function gcd_lcm(a, b) {
  function gcd(x, y) {
    while (y) {
      [x, y] = [y, x % y];
    }
    return x;
  }
  const g = gcd(a, b);
  const l = Math.floor(a * b / g);
  return [g, l];
}
