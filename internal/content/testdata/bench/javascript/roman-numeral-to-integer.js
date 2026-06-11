function roman_to_int(s) {
    const doubles = {'CM': 900, 'CD': 400, 'XC': 90, 'XL': 40, 'IX': 9, 'IV': 4};
    const singles = {'M': 1000, 'D': 500, 'C': 100, 'L': 50, 'X': 10, 'V': 5, 'I': 1};
    let total = 0;
    let i = 0;
    while (i < s.length) {
        if (i < s.length - 1 && s.slice(i, i+2) in doubles) {
            total += doubles[s.slice(i, i+2)];
            i += 2;
        } else {
            total += singles[s[i]];
            i += 1;
        }
    }
    return total;
}
