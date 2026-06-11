function eval_rpn(tokens: string[]): number {
    const stack: number[] = [];
    const ops = new Set(['+', '-', '*', '/']);
    for (const t of tokens) {
        if (ops.has(t)) {
            const b = stack.pop()!;
            const a = stack.pop()!;
            if (t === '+') stack.push(a + b);
            else if (t === '-') stack.push(a - b);
            else if (t === '*') stack.push(a * b);
            else stack.push(Math.trunc(a / b));
        } else {
            stack.push(parseInt(t, 10));
        }
    }
    return stack[0];
}
