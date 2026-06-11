function count_pairs_same_mod(nums: number[], divisor: number): number {
    const freq: {[key: number]: number} = {};
    for (const x of nums) {
        const r = x % divisor;
        freq[r] = (freq[r] || 0) + 1;
    }
    let pairs = 0;
    for (const cnt of Object.values(freq)) {
        pairs += cnt * (cnt - 1) / 2;
    }
    return pairs;
}
