function check_valid_string(s) {
    let lo = 0, hi = 0;
    for (const c of s) {
        if (c === '(') { lo++; hi++; }
        else if (c === ')') { lo--; hi--; }
        else { lo--; hi++; }
        if (hi < 0) return false;
        if (lo < 0) lo = 0;
    }
    return lo === 0;
}
