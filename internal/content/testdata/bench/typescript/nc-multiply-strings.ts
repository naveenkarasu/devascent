function multiply(num1: string, num2: string): string {
    if (num1 === "0" || num2 === "0") return "0";
    const m = num1.length, n = num2.length;
    const buf = new Array(m + n).fill(0);
    for (let i = m - 1; i >= 0; i--) {
        for (let j = n - 1; j >= 0; j--) {
            const mul = (num1.charCodeAt(i) - 48) * (num2.charCodeAt(j) - 48);
            const p1 = i + j, p2 = i + j + 1;
            const total = mul + buf[p2];
            buf[p2] = total % 10;
            buf[p1] += Math.floor(total / 10);
        }
    }
    const result = buf.join('').replace(/^0+/, '');
    return result || "0";
}
