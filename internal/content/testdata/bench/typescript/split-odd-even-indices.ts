function split_odd_even(s: string): string[] {
    let even_chars = '';
    let odd_chars = '';
    for (let i = 0; i < s.length; i++) {
        if (i % 2 === 0) even_chars += s[i];
        else odd_chars += s[i];
    }
    return [even_chars, odd_chars];
}
