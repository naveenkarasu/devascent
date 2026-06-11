function reverse_integer(x) {
    const INT_MIN = -(2 ** 31), INT_MAX = 2 ** 31 - 1;
    const sign = x < 0 ? -1 : 1;
    const digits = Math.abs(x).toString().split('').reverse().join('');
    const result = sign * parseInt(digits, 10);
    if (result < INT_MIN || result > INT_MAX) return 0;
    return result;
}
