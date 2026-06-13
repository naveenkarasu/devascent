function is_perfect_square(n: number): boolean {
    if (n < 0) return false;
    let r = Math.floor(Math.sqrt(n));
    // Adjust for floating-point inaccuracy
    while (r * r > n) r--;
    while ((r + 1) * (r + 1) <= n) r++;
    return r * r === n;
}
