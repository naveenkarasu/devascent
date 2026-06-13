function ladder_length(begin_word: string, end_word: string, word_list: string[]): number {
    const wordSet = new Set(word_list);
    if (!wordSet.has(end_word)) return 0;
    const queue: [string, number][] = [[begin_word, 1]];
    const visited = new Set<string>([begin_word]);
    let qi = 0;
    while (qi < queue.length) {
        const [word, length] = queue[qi++];
        for (let i = 0; i < word.length; i++) {
            for (let c = 97; c <= 122; c++) {
                const candidate = word.slice(0, i) + String.fromCharCode(c) + word.slice(i + 1);
                if (candidate === end_word) return length + 1;
                if (wordSet.has(candidate) && !visited.has(candidate)) {
                    visited.add(candidate);
                    queue.push([candidate, length + 1]);
                }
            }
        }
    }
    return 0;
}
