function find_unique_among_k(k, arr) {
    let counts = new Map();
    for (let x of arr) {
        counts.set(x, (counts.get(x) || 0) + 1);
    }
    for (let [val, cnt] of counts) {
        if (cnt % k !== 0) return val;
    }
    return -1;
}
