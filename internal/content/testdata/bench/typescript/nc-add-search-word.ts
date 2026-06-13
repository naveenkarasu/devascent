function word_dictionary_ops(operations: any[][]): any[] {
    const root: any = {};

    function add_word(word: string): void {
        let node = root;
        for (const c of word) {
            if (!node[c]) node[c] = {};
            node = node[c];
        }
        node['$'] = true;
    }

    function search(word: string): boolean {
        function dfs(node: any, i: number): boolean {
            if (i === word.length) return '$' in node;
            const c = word[i];
            if (c === '.') {
                for (const k of Object.keys(node)) {
                    if (k !== '$' && dfs(node[k], i + 1)) return true;
                }
                return false;
            }
            if (!(c in node)) return false;
            return dfs(node[c], i + 1);
        }
        return dfs(root, 0);
    }

    const out: any[] = [];
    for (const [op, arg] of operations) {
        if (op === 'addWord') {
            add_word(arg);
            out.push(null);
        } else {
            out.push(search(arg));
        }
    }
    return out;
}
