function detect_squares_ops(operations) {
    const cnt = new Map(); // "x,y" -> count

    function key(x, y) { return x + ',' + y; }
    function getCount(x, y) { return cnt.get(key(x, y)) || 0; }

    const out = [];
    for (const op of operations) {
        if (op[0] === "add") {
            const [, pt] = op;
            const k = key(pt[0], pt[1]);
            cnt.set(k, (cnt.get(k) || 0) + 1);
            out.push(null);
        } else {
            // count
            const [, pt] = op;
            const px = pt[0], py = pt[1];
            let total = 0;
            for (const [k, c] of cnt) {
                const [xs, ys] = k.split(',');
                const x = parseInt(xs), y = parseInt(ys);
                if (Math.abs(x - px) === Math.abs(y - py) && x !== px && y !== py) {
                    total += c * getCount(px, y) * getCount(x, py);
                }
            }
            out.push(total);
        }
    }
    return out;
}
