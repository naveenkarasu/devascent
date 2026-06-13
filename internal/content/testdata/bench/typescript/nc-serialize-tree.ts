function serialize(root: any): string {
    const out: string[] = [];
    function dfs(node: any): void {
        if (node === null) {
            out.push('#');
            return;
        }
        out.push(String(node.val));
        dfs(node.left);
        dfs(node.right);
    }
    dfs(root);
    return out.join(',');
}

function deserialize(data: string): any {
    const vals = data.split(',');
    let idx = 0;
    function build(): any {
        const v = vals[idx++];
        if (v === '#') return null;
        const node = { val: parseInt(v), left: null, right: null };
        node.left = build();
        node.right = build();
        return node;
    }
    return build();
}

function codec_roundtrip(root: any): any {
    const data = serialize(root);
    return deserialize(data);
}
