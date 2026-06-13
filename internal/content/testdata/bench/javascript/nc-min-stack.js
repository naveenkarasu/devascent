function min_stack_ops(operations) {
    const stack = [];
    const mins = [];
    const out = [];
    for (const op of operations) {
        if (op[0] === "push") {
            stack.push(op[1]);
            mins.push(mins.length === 0 ? op[1] : Math.min(op[1], mins[mins.length - 1]));
            out.push(null);
        } else if (op[0] === "pop") {
            stack.pop();
            mins.pop();
            out.push(null);
        } else if (op[0] === "top") {
            out.push(stack[stack.length - 1]);
        } else {
            out.push(mins[mins.length - 1]);
        }
    }
    return out;
}
