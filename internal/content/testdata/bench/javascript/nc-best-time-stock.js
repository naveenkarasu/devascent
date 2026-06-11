function max_profit(prices) {
    let minPrice = Infinity;
    let best = 0;
    for (const p of prices) {
        if (p < minPrice) {
            minPrice = p;
        } else if (p - minPrice > best) {
            best = p - minPrice;
        }
    }
    return best;
}
