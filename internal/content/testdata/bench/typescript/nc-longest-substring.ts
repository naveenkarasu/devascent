function length_of_longest_substring(s: string): number {
    const lastSeen: {[ch: string]: number} = {};
    let left = 0;
    let best = 0;
    for (let right = 0; right < s.length; right++) {
        const ch = s[right];
        if (ch in lastSeen && lastSeen[ch] >= left) {
            left = lastSeen[ch] + 1;
        }
        lastSeen[ch] = right;
        if (right - left + 1 > best) {
            best = right - left + 1;
        }
    }
    return best;
}
