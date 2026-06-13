function letter_combinations(digits) {
    if (!digits) return [];
    const mapping = {
        '2': 'abc', '3': 'def', '4': 'ghi', '5': 'jkl',
        '6': 'mno', '7': 'pqrs', '8': 'tuv', '9': 'wxyz'
    };
    let results = [''];
    for (const d of digits) {
        const newResults = [];
        for (const prev of results) {
            for (const ch of mapping[d]) {
                newResults.push(prev + ch);
            }
        }
        results = newResults;
    }
    return results;
}
