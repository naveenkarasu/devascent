function merge_k_lists(lists) {
    const all = [];
    for (const lst of lists) {
        for (const v of lst) {
            all.push(v);
        }
    }
    all.sort((a, b) => a - b);
    return all;
}
