function sequence_equation(p) {
    let n = p.length;
    let inv = new Array(n + 1).fill(0);
    for (let i = 0; i < n; i++) {
        inv[p[i]] = i + 1;
    }
    let result = [];
    for (let x = 1; x <= n; x++) {
        result.push(inv[inv[x]]);
    }
    return result;
}
