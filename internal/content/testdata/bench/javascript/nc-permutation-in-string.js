function check_inclusion(s1, s2) {
    if (s1.length > s2.length) return false;
    const need = new Map();
    for (const ch of s1) {
        need.set(ch, (need.get(ch) || 0) + 1);
    }
    const window = new Map();
    const k = s1.length;
    for (let i = 0; i < s2.length; i++) {
        const ch = s2[i];
        window.set(ch, (window.get(ch) || 0) + 1);
        if (i >= k) {
            const lch = s2[i - k];
            window.set(lch, window.get(lch) - 1);
            if (window.get(lch) === 0) {
                window.delete(lch);
            }
        }
        // Compare maps
        if (window.size === need.size) {
            let match = true;
            for (const [key, val] of need) {
                if (window.get(key) !== val) {
                    match = false;
                    break;
                }
            }
            if (match) return true;
        }
    }
    return false;
}
