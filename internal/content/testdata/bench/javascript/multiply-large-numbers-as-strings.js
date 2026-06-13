function multiply_strings(num1, num2) {
    const m = num1.length, n = num2.length;
    const pos = new Array(m + n).fill(0);
    for (let i = m - 1; i >= 0; i--) {
        for (let j = n - 1; j >= 0; j--) {
            const mul = (num1.charCodeAt(i) - 48) * (num2.charCodeAt(j) - 48);
            const p1 = i + j, p2 = i + j + 1;
            const total = mul + pos[p2];
            pos[p2] = total % 10;
            pos[p1] += Math.floor(total / 10);
        }
    }
    const s = pos.join('').replace(/^0+/, '');
    return s || '0';
}
