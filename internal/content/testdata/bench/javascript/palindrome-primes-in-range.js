function palindrome_primes_in_range(a, b) {
  function isPrime(n) {
    if (n < 2) return false;
    let r = Math.floor(Math.sqrt(n));
    while (r * r > n) r--;
    while ((r + 1) * (r + 1) <= n) r++;
    for (let i = 2; i <= r; i++) {
      if (n % i === 0) return false;
    }
    return true;
  }
  function isPalindrome(n) {
    const s = String(n);
    return s === s.split('').reverse().join('');
  }
  const result = [];
  for (let i = a; i <= b; i++) {
    if (isPrime(i) && isPalindrome(i)) result.push(i);
  }
  return result;
}
