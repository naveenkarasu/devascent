function sum_of_factorials(n) {
  let total = BigInt(0);
  let fact = BigInt(1);
  for (let i = 1; i <= n; i++) {
    fact *= BigInt(i);
    total += fact;
  }
  return Number(total);
}
