fn divide_integers(dividend: i64, divisor: i64) -> i64 {
    let int_max: i64 = 2147483647; // 2^31 - 1
    let int_min: i64 = -2147483648; // -2^31
    if dividend == int_min && divisor == -1 {
        return int_max;
    }
    let negative = (dividend < 0) != (divisor < 0);
    let mut a = dividend.abs();
    let b = divisor.abs();
    let mut result: i64 = 0;
    while a >= b {
        let mut temp = b;
        let mut multiple: i64 = 1;
        while a >= (temp << 1) {
            temp <<= 1;
            multiple <<= 1;
        }
        a -= temp;
        result += multiple;
    }
    if negative {
        -result
    } else {
        result
    }
}
