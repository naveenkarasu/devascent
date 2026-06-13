function has_adjacent_zeros(s: string): boolean {
    for (let i = 1; i < s.length; i++) {
        if (s[i] === '0' && s[i - 1] === '0') return true;
    }
    return false;
}
