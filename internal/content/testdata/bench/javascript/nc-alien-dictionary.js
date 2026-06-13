function alien_order(words) {
    const adj = new Map();
    const indeg = new Map();

    // Initialize all chars
    for (const w of words) {
        for (const c of w) {
            if (!indeg.has(c)) indeg.set(c, 0);
            if (!adj.has(c)) adj.set(c, new Set());
        }
    }

    // Build graph from word pairs
    for (let i = 0; i < words.length - 1; i++) {
        const a = words[i], b = words[i + 1];
        const m = Math.min(a.length, b.length);
        if (a.length > b.length && a.slice(0, m) === b.slice(0, m)) return "";
        for (let j = 0; j < m; j++) {
            if (a[j] !== b[j]) {
                if (!adj.get(a[j]).has(b[j])) {
                    adj.get(a[j]).add(b[j]);
                    indeg.set(b[j], indeg.get(b[j]) + 1);
                }
                break;
            }
        }
    }

    // BFS with always picking smallest char
    const queue = Array.from(indeg.keys()).filter(c => indeg.get(c) === 0).sort();
    const res = [];
    while (queue.length > 0) {
        queue.sort();
        const c = queue.shift();
        res.push(c);
        const neighbors = Array.from(adj.get(c)).sort();
        for (const nxt of neighbors) {
            indeg.set(nxt, indeg.get(nxt) - 1);
            if (indeg.get(nxt) === 0) {
                queue.push(nxt);
            }
        }
    }

    if (res.length !== indeg.size) return "";
    return res.join('');
}
