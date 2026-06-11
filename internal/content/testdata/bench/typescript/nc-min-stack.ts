class MinStack {
    private stack: number[];
    private mins: number[];
    constructor() {
        this.stack = [];
        this.mins = [];
    }
    push(val: number): void {
        this.stack.push(val);
        this.mins.push(this.mins.length === 0 ? val : Math.min(val, this.mins[this.mins.length - 1]));
    }
    pop(): void {
        this.stack.pop();
        this.mins.pop();
    }
    top(): number {
        return this.stack[this.stack.length - 1];
    }
    get_min(): number {
        return this.mins[this.mins.length - 1];
    }
}

function min_stack_ops(operations: any[][]): any[] {
    const st = new MinStack();
    const out: any[] = [];
    for (const op of operations) {
        if (op[0] === "push") {
            st.push(op[1]);
            out.push(null);
        } else if (op[0] === "pop") {
            st.pop();
            out.push(null);
        } else if (op[0] === "top") {
            out.push(st.top());
        } else {
            out.push(st.get_min());
        }
    }
    return out;
}
