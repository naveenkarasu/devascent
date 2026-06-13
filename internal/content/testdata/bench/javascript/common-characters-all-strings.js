function count_common_characters(strings) {
    if (!strings || strings.length === 0) return 0;
    let common = new Set(strings[0]);
    for (let i = 1; i < strings.length; i++) {
        let s = new Set(strings[i]);
        for (let ch of common) {
            if (!s.has(ch)) common.delete(ch);
        }
    }
    return common.size;
}
