function word_break(s: string, word_dict: string[]): boolean {
    const word_set = new Set(word_dict);
    const n = s.length;
    const dp: boolean[] = new Array(n + 1).fill(false);
    dp[0] = true;
    for (let i = 1; i <= n; i++) {
        for (let j = 0; j < i; j++) {
            if (dp[j] && word_set.has(s.slice(j, i))) {
                dp[i] = true;
                break;
            }
        }
    }
    return dp[n];
}
