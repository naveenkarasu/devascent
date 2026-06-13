function partition(s) {
    const res = [];
    function is_palindrome(t) {
        return t === t.split('').reverse().join('');
    }
    function backtrack(start, current) {
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
    res.sort((a, b) => {
        for (let i = 0; i < Math.min(a.length, b.length); i++) {
            if (a[i] < b[i]) return -1;
            if (a[i] > b[i]) return 1;
        }
        return a.length - b.length;
    });
    return res;
}
