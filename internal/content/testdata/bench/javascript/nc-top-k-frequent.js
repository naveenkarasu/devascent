function top_k_frequent(nums, k) {
    const count = new Map();
    for (const n of nums) {
        count.set(n, (count.get(n) || 0) + 1);
    }
    const keys = Array.from(count.keys());
    keys.sort((a, b) => {
        const diff = count.get(b) - count.get(a);
        if (diff !== 0) return diff;
        return a - b;
    });
    const top = keys.slice(0, k);
    top.sort((a, b) => a - b);
    return top;
}
