function partition(s: string): string[][] {
    function cmpArr(a: string[], b: string[]): number {
        const n = Math.min(a.length, b.length);
        for (let i = 0; i < n; i++) {
            if (a[i] < b[i]) return -1;
            if (a[i] > b[i]) return 1;
        }
        return a.length - b.length;
    }

    const res: string[][] = [];

    function is_palindrome(t: string): boolean {
        return t === t.split('').reverse().join('');
    }

    function backtrack(start: number, current: string[]): void {
        if (start === s.length) {
            res.push([...current]);
            return;
        }
        for (let end = start + 1; end <= s.length; end++) {
            const substr = s.slice(start, end);
            if (is_palindrome(substr)) {
                current.push(substr);
                backtrack(end, current);
                current.pop();
            }
        }
    }

    backtrack(0, []);
    res.sort(cmpArr);
    return res;
}
