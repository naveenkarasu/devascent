function is_valid(s: string): boolean {
    const stack: string[] = [];
    const pairs: {[ch: string]: string} = {')': '(', ']': '[', '}': '{'};
    for (const ch of s) {
        if (ch in pairs) {
            if (stack.length === 0 || stack[stack.length - 1] !== pairs[ch]) {
                return false;
            }
            stack.pop();
        } else {
            stack.push(ch);
        }
    }
    return stack.length === 0;
}
