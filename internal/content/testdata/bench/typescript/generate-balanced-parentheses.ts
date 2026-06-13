function generate_parentheses(n: number): string[] {
    const results: string[] = [];
    function backtrack(s: string, open_count: number, close_count: number): void {
        if (s.length === 2 * n) {
            results.push(s);
            return;
        }
        if (open_count < n) {
            backtrack(s + '(', open_count + 1, close_count);
        }
        if (close_count < open_count) {
            backtrack(s + ')', open_count, close_count + 1);
        }
    }
    backtrack('', 0, 0);
    return results;
}
