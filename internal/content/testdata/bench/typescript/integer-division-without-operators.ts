function divide_integers(dividend: number, divisor: number): number {
    const INT_MAX = 2147483647;
    const INT_MIN = -2147483648;
    if (dividend === INT_MIN && divisor === -1) return INT_MAX;
    const negative = (dividend < 0) !== (divisor < 0);
    let a = Math.abs(dividend);
    let b = Math.abs(divisor);
    let result = 0;
    while (a >= b) {
        let temp = b, multiple = 1;
        while (a >= (temp * 2)) {
            temp *= 2;
            multiple *= 2;
        }
        a -= temp;
        result += multiple;
    }
    return negative ? -result : result;
}
