class DetectSquares {
    private cnt: Map<string, number>;
    constructor() {
        this.cnt = new Map();
    }
    add(point: number[]): void {
        const key = `${point[0]},${point[1]}`;
        this.cnt.set(key, (this.cnt.get(key) || 0) + 1);
    }
    count(point: number[]): number {
        const [px, py] = point;
        let total = 0;
        for (const [k, c] of this.cnt) {
            const [xs, ys] = k.split(',');
            const x = parseInt(xs), y = parseInt(ys);
            if (Math.abs(x - px) === Math.abs(y - py) && x !== px && y !== py) {
                total += c * (this.cnt.get(`${px},${y}`) || 0) * (this.cnt.get(`${x},${py}`) || 0);
            }
        }
        return total;
    }
}

function detect_squares_ops(operations: any[][]): any[] {
    const ds = new DetectSquares();
    const out: any[] = [];
    for (const op of operations) {
        if (op[0] === "add") {
            ds.add(op[1]);
            out.push(null);
        } else {
            out.push(ds.count(op[1]));
        }
    }
    return out;
}
