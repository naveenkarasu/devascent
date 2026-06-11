fn get_sum(a: i64, b: i64) -> i64 {
    let mask: i64 = 0xFFFFFFFF;
    let mut a = a & mask;
    let mut b = b & mask;
    while (b & mask) != 0 {
        let carry = ((a & b) << 1) & mask;
        a = (a ^ b) & mask;
        b = carry;
    }
    a &= mask;
    if a <= 0x7FFFFFFF {
        a
    } else {
        !(a ^ mask)
    }
}
