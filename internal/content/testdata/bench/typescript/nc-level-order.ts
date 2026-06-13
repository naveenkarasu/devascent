function level_order(root: any): number[][] {
    if (root === null) return [];
    const q: any[] = [root];
    const res: number[][] = [];
    while (q.length > 0) {
        const level: number[] = [];
        const size = q.length;
        for (let i = 0; i < size; i++) {
            const node = q.shift();
            level.push(node.val);
            if (node.left !== null) q.push(node.left);
            if (node.right !== null) q.push(node.right);
        }
        res.push(level);
    }
    return res;
}
