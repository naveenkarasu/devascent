function nth_prime(n) {
  const primes = [];
  let i = 2;
  while (primes.length < n) {
    if (primes.every(p => i % p !== 0)) {
      primes.push(i);
    }
    i++;
  }
  return primes[primes.length - 1];
}
