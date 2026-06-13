function find_words(board, words) {
    // Build trie
    const trie = {};
    for (const w of words) {
        let node = trie;
        for (const c of w) {
            if (!node[c]) node[c] = {};
            node = node[c];
        }
        node['$'] = w;
    }

    if (!board || board.length === 0) return [];
    const rows = board.length, cols = board[0].length;
    const found = new Set();

    function dfs(r, c, node) {
        const ch = board[r][c];
        if (!(ch in node)) return;
        const nxt = node[ch];
        if ('$' in nxt) found.add(nxt['$']);
        board[r][c] = '#';
        const dirs = [[1, 0], [-1, 0], [0, 1], [0, -1]];
        for (const [dr, dc] of dirs) {
            const nr = r + dr, nc = c + dc;
            if (nr >= 0 && nr < rows && nc >= 0 && nc < cols && board[nr][nc] !== '#') {
                dfs(nr, nc, nxt);
            }
        }
        board[r][c] = ch;
    }

    for (let r = 0; r < rows; r++) {
        for (let c = 0; c < cols; c++) {
            dfs(r, c, trie);
        }
    }
    return Array.from(found).sort();
}
