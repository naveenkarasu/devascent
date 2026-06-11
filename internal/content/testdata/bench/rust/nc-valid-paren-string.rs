fn check_valid_string(s: String) -> bool {
    let mut lo = 0i64;
    let mut hi = 0i64;
    for c in s.chars() {
        if c == '(' {
            lo += 1;
            hi += 1;
        } else if c == ')' {
            lo -= 1;
            hi -= 1;
        } else {
            lo -= 1;
            hi += 1;
        }
        if hi < 0 {
            return false;
        }
        if lo < 0 {
            lo = 0;
        }
    }
    lo == 0
}
