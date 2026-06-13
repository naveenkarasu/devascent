function merge_k_lists(lists: number[][]): number[] {
    const out: number[] = [];
    for (const lst of lists) {
        for (const v of lst) {
            out.push(v);
        }
    }
    out.sort((a, b) => a - b);
    return out;
}
