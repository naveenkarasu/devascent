function generate_parentheses(n: number): string[] {
    const result: string[] = [];
    function backtrack(s: string, openCount: number, closeCount: number): void {
        if (s.length === 2 * n) {
            result.push(s);
            return;
        }
        if (openCount < n) backtrack(s + '(', openCount + 1, closeCount);
        if (closeCount < openCount) backtrack(s + ')', openCount, closeCount + 1);
    }
    backtrack('', 0, 0);
    return result.sort();
}
