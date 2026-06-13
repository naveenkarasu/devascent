function top_k_frequent(nums: number[], k: number): number[] {
    const cnt: Record<number, number> = {};
    for (const n of nums) {
        cnt[n] = (cnt[n] || 0) + 1;
    }
    const keys = Object.keys(cnt).map(Number);
    keys.sort((a, b) => cnt[b] - cnt[a] || a - b);
    const top = keys.slice(0, k);
    top.sort((a, b) => a - b);
    return top;
}
