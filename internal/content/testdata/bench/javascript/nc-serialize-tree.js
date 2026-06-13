function codec_roundtrip(root) {
    // Serialize
    const out = [];
    function dfs(node) {
        if (node === null) {
            out.push('#');
            return;
        }
        out.push(String(node.val));
        dfs(node.left);
        dfs(node.right);
    }
    dfs(root);
    const data = out.join(',');

    // Deserialize
    const vals = data.split(',');
    let idx = 0;
    function build() {
        const v = vals[idx++];
        if (v === '#') return null;
        const node = { val: parseInt(v), left: null, right: null };
        node.left = build();
        node.right = build();
        return node;
    }
    return build();
}
