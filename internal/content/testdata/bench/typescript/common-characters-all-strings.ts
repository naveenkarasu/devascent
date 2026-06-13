function count_common_characters(strings: string[]): number {
    if (strings.length === 0) return 0;
    let common = new Set<string>(strings[0].split(''));
    for (let i = 1; i < strings.length; i++) {
        const s = new Set<string>(strings[i].split(''));
        for (const ch of common) {
            if (!s.has(ch)) common.delete(ch);
        }
    }
    return common.size;
}
