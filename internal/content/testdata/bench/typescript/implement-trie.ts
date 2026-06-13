function trie_ops(operations: any[][]): any[] {
    const root: any = {};

    function insert(word: string): null {
        let node = root;
        for (const c of word) {
            if (!(c in node)) node[c] = {};
            node = node[c];
        }
        node["$"] = true;
        return null;
    }

    function search(word: string): boolean {
        let node = root;
        for (const c of word) {
            if (!(c in node)) return false;
            node = node[c];
        }
        return "$" in node;
    }

    function starts_with(prefix: string): boolean {
        let node = root;
        for (const c of prefix) {
            if (!(c in node)) return false;
            node = node[c];
        }
        return true;
    }

    const out: any[] = [];
    for (const [op, arg] of operations) {
        if (op === "insert") {
            out.push(insert(arg));
        } else if (op === "search") {
            out.push(search(arg));
        } else {
            out.push(starts_with(arg));
        }
    }
    return out;
}
