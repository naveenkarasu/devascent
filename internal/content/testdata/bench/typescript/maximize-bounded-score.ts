function maximize_score(arr: number[], k: number): number {
    let best = 0;
    for (const x of arr) {
        const score = Math.floor((2 * k) / (1 + 2 * Math.abs(x - k)));
        if (score > best) best = score;
    }
    return best;
}
