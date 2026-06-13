function max_square_divisor_root(n: number): number {
    let answer = 1;
    let remaining = n;
    let fact = 2;
    while (fact * fact <= n) {
        let count = 0;
        while (remaining % fact === 0) {
            count++;
            remaining = Math.floor(remaining / fact);
            if (count === 2) {
                answer *= fact;
                count = 0;
            }
        }
        fact++;
    }
    return answer;
}
