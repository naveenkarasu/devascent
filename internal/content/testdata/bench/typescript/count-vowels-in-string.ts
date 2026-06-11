function count_vowels(s: string): number {
    const vowels = new Set(['a','e','i','o','u','A','E','I','O','U']);
    let count = 0;
    for (const ch of s) {
        if (vowels.has(ch)) count++;
    }
    return count;
}
