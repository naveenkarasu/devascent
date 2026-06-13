function right_side_view(root: any): number[] {
    if (root === null) return [];
    const q: any[] = [root];
    const res: number[] = [];
    while (q.length > 0) {
        const n = q.length;
        for (let i = 0; i < n; i++) {
            const node = q.shift();
            if (i === n - 1) res.push(node.val);
            if (node.left !== null) q.push(node.left);
            if (node.right !== null) q.push(node.right);
        }
    }
    return res;
}
