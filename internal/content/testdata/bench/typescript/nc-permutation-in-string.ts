function check_inclusion(s1: string, s2: string): boolean {
    if (s1.length > s2.length) return false;
    const need: {[ch: string]: number} = {};
    for (const ch of s1) {
        need[ch] = (need[ch] || 0) + 1;
    }
    const window: {[ch: string]: number} = {};
    const k = s1.length;
    for (let i = 0; i < s2.length; i++) {
        const ch = s2[i];
        window[ch] = (window[ch] || 0) + 1;
        if (i >= k) {
            const leftCh = s2[i - k];
            window[leftCh] -= 1;
            if (window[leftCh] === 0) delete window[leftCh];
        }
        // Compare window and need
        const needKeys = Object.keys(need);
        const windowKeys = Object.keys(window);
        if (needKeys.length === windowKeys.length &&
            needKeys.every(c => need[c] === window[c])) {
            return true;
        }
    }
    return false;
}
