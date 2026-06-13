function compressed_length(s: string): number {
    const n = s.length;
    let count = 0;
    let i = 0;
    while (i < n) {
        if (i + 1 < n && s[i] === s[i + 1]) {
            count++;
            i += 2;
        } else {
            count++;
            i++;
        }
    }
    return count;
}
