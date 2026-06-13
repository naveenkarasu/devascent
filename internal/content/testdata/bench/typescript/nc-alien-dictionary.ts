function alien_order(words: string[]): string {
    const adj: Record<string, Set<string>> = {};
    const indeg: Record<string, number> = {};

    for (const w of words) {
        for (const c of w) {
            if (!(c in indeg)) indeg[c] = 0;
            if (!(c in adj)) adj[c] = new Set();
        }
    }

    for (let i = 0; i < words.length - 1; i++) {
        const a = words[i], b = words[i + 1];
        const m = Math.min(a.length, b.length);
        if (a.length > b.length && a.slice(0, m) === b.slice(0, m)) return "";
        for (let j = 0; j < m; j++) {
            if (a[j] !== b[j]) {
                if (!adj[a[j]].has(b[j])) {
                    adj[a[j]].add(b[j]);
                    indeg[b[j]]++;
                }
                break;
            }
        }
    }

    // min-heap via sorted array (alphabet is small)
    let queue = Object.keys(indeg).filter(c => indeg[c] === 0).sort();
    const res: string[] = [];
    while (queue.length > 0) {
        const c = queue.shift()!;
        res.push(c);
        const neighbors = Array.from(adj[c]).sort();
        for (const nxt of neighbors) {
            indeg[nxt]--;
            if (indeg[nxt] === 0) {
                queue.push(nxt);
                queue.sort();
            }
        }
    }

    if (res.length !== Object.keys(indeg).length) return "";
    return res.join('');
}
