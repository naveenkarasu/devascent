function min_window(s: string, t: string): string {
    if (!t || !s) return "";
    const need: {[ch: string]: number} = {};
    for (const ch of t) {
        need[ch] = (need[ch] || 0) + 1;
    }
    const have: {[ch: string]: number} = {};
    let formed = 0;
    const required = Object.keys(need).length;
    let left = 0;
    let bestLen = Infinity;
    let bestStart = 0;
    for (let right = 0; right < s.length; right++) {
        const ch = s[right];
        have[ch] = (have[ch] || 0) + 1;
        if (ch in need && have[ch] === need[ch]) formed++;
        while (formed === required) {
            if (right - left + 1 < bestLen) {
                bestLen = right - left + 1;
                bestStart = left;
            }
            const leftCh = s[left];
            have[leftCh] -= 1;
            if (leftCh in need && have[leftCh] < need[leftCh]) formed--;
            left++;
        }
    }
    if (bestLen === Infinity) return "";
    return s.slice(bestStart, bestStart + bestLen);
}
