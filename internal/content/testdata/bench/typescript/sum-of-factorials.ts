function sum_of_factorials(n: number): number {
    let total = 0;
    let fact = 1;
    for (let i = 1; i <= n; i++) {
        fact *= i;
        total += fact;
    }
    return total;
}
