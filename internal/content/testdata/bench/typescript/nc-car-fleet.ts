function car_fleet(target: number, position: number[], speed: number[]): number {
    const pairs: [number, number][] = position.map((p, i) => [p, speed[i]]);
    pairs.sort((a, b) => b[0] - a[0]);
    const stack: number[] = [];
    for (const [pos, spd] of pairs) {
        const time = (target - pos) / spd;
        if (stack.length === 0 || time > stack[stack.length - 1]) {
            stack.push(time);
        }
    }
    return stack.length;
}
