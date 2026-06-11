function product_of_counts(s: string): number {
    const MOD = 1000000007;
    const freq: {[key: string]: number} = {};
    for (const ch of s) {
        freq[ch] = (freq[ch] || 0) + 1;
    }
    let result = 1;
    for (const cnt of Object.values(freq)) {
        result = (result * cnt) % MOD;
    }
    return result;
}
