function divide_integers(dividend, divisor) {
    const INT_MAX = 2147483647;
    const INT_MIN = -2147483648;
    if (dividend === INT_MIN && divisor === -1) return INT_MAX;
    const negative = (dividend < 0) !== (divisor < 0);
    let a = Math.abs(dividend), b = Math.abs(divisor);
    let result = 0;
    while (a >= b) {
        let temp = b, multiple = 1;
        while (a >= (temp << 1)) {
            temp <<= 1;
            multiple <<= 1;
        }
        a -= temp;
        result += multiple;
    }
    return negative ? -result : result;
}
