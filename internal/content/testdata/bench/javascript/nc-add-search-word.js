function word_dictionary_ops(operations) {
    const root = {};

    function addWord(word) {
        let node = root;
        for (const c of word) {
            if (!node[c]) node[c] = {};
            node = node[c];
        }
        node['$'] = true;
    }

    function search(word) {
        function dfs(node, i) {
            if (i === word.length) return '$' in node;
            const c = word[i];
            if (c === '.') {
                for (const k in node) {
                    if (k !== '$' && dfs(node[k], i + 1)) return true;
                }
                return false;
            }
            if (!(c in node)) return false;
            return dfs(node[c], i + 1);
        }
        return dfs(root, 0);
    }

    const out = [];
    for (const op of operations) {
        const opName = op[0], arg = op[1];
        if (opName === 'addWord') {
            addWord(arg);
            out.push(null);
        } else {
            out.push(search(arg));
        }
    }
    return out;
}
