function car_fleet(target, position, speed) {
    const n = position.length;
    const pairs = [];
    for (let i = 0; i < n; i++) {
        pairs.push([position[i], speed[i]]);
    }
    pairs.sort((a, b) => b[0] - a[0]);
    const stack = [];
    for (const [pos, spd] of pairs) {
        const time = (target - pos) / spd;
        if (stack.length === 0 || time > stack[stack.length - 1]) {
            stack.push(time);
        }
    }
    return stack.length;
}
