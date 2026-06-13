function split_odd_even(s) {
    let even = '';
    let odd = '';
    for (let i = 0; i < s.length; i++) {
        if (i % 2 === 0) even += s[i];
        else odd += s[i];
    }
    return [even, odd];
}
