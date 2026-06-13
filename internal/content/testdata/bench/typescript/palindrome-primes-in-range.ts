function palindrome_primes_in_range(a: number, b: number): number[] {
    function isPrime(n: number): boolean {
        if (n < 2) return false;
        for (let i = 2; i * i <= n; i++) {
            if (n % i === 0) return false;
        }
        return true;
    }
    function isPalindrome(n: number): boolean {
        const s = String(n);
        return s === s.split('').reverse().join('');
    }
    const result: number[] = [];
    for (let i = a; i <= b; i++) {
        if (isPrime(i) && isPalindrome(i)) result.push(i);
    }
    return result;
}
