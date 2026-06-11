function encode(strs: string[]): string {
    return strs.map(s => s.length + '#' + s).join('');
}

function decode(s: string): string[] {
    const res: string[] = [];
    let i = 0;
    while (i < s.length) {
        let j = i;
        while (s[j] !== '#') j++;
        const length = parseInt(s.slice(i, j));
        res.push(s.slice(j + 1, j + 1 + length));
        i = j + 1 + length;
    }
    return res;
}

function encode_decode(strs: string[]): string[] {
    const s = encode(strs);
    return decode(s);
}
