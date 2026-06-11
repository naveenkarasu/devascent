function min_window(s, t) {
    if (!t || !s) return "";
    const need = new Map();
    for (const ch of t) {
        need.set(ch, (need.get(ch) || 0) + 1);
    }
    const have = new Map();
    let formed = 0;
    const required = need.size;
    let left = 0;
    let bestLen = Infinity;
    let bestStart = 0;
    for (let right = 0; right < s.length; right++) {
        const ch = s[right];
        have.set(ch, (have.get(ch) || 0) + 1);
        if (need.has(ch) && have.get(ch) === need.get(ch)) {
            formed++;
        }
        while (formed === required) {
            if (right - left + 1 < bestLen) {
                bestLen = right - left + 1;
                bestStart = left;
            }
            const lch = s[left];
            have.set(lch, have.get(lch) - 1);
            if (need.has(lch) && have.get(lch) < need.get(lch)) {
                formed--;
            }
            left++;
        }
    }
    if (bestLen === Infinity) return "";
    return s.substring(bestStart, bestStart + bestLen);
}
