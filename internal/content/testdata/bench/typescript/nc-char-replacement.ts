function character_replacement(s: string, k: number): number {
    const counts: {[ch: string]: number} = {};
    let left = 0;
    let maxCount = 0;
    let best = 0;
    for (let right = 0; right < s.length; right++) {
        const ch = s[right];
        counts[ch] = (counts[ch] || 0) + 1;
        if (counts[ch] > maxCount) maxCount = counts[ch];
        while ((right - left + 1) - maxCount > k) {
            counts[s[left]] -= 1;
            left++;
        }
        if (right - left + 1 > best) best = right - left + 1;
    }
    return best;
}
