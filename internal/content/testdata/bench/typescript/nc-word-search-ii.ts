function find_words(board: string[][], words: string[]): string[] {
    const trie: any = {};
    for (const w of words) {
        let node = trie;
        for (const c of w) {
            if (!node[c]) node[c] = {};
            node = node[c];
        }
        node['$'] = w;
    }

    if (!board || board.length === 0 || board[0].length === 0) return [];
    const rows = board.length, cols = board[0].length;
    const found = new Set<string>();

    function dfs(r: number, c: number, node: any): void {
        const ch = board[r][c];
        if (!(ch in node)) return;
        const nxt = node[ch];
        if ('$' in nxt) found.add(nxt['$']);
        board[r][c] = '#';
        for (const [dr, dc] of [[1,0],[-1,0],[0,1],[0,-1]]) {
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
