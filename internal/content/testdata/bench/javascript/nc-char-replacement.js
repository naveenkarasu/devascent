function character_replacement(s, k) {
    const counts = new Map();
    let left = 0;
    let maxCount = 0;
    let best = 0;
    for (let right = 0; right < s.length; right++) {
        const ch = s[right];
        counts.set(ch, (counts.get(ch) || 0) + 1);
        if (counts.get(ch) > maxCount) {
            maxCount = counts.get(ch);
        }
        while ((right - left + 1) - maxCount > k) {
            const lch = s[left];
            counts.set(lch, counts.get(lch) - 1);
            left++;
        }
        if (right - left + 1 > best) {
            best = right - left + 1;
        }
    }
    return best;
}
