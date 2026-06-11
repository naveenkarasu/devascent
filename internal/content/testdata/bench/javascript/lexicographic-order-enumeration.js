function lexical_order(n) {
    const result = [];
    let cur = 1;
    while (result.length < n) {
        result.push(cur);
        if (cur * 10 <= n) {
            cur *= 10;
        } else {
            while (cur % 10 === 9 || cur >= n) {
                cur = Math.floor(cur / 10);
            }
            cur += 1;
        }
    }
    return result;
}
