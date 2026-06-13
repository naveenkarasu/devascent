function daily_temperatures(temps: number[]): number[] {
    const result = new Array(temps.length).fill(0);
    const stack: [number, number][] = []; // [temperature, index]
    for (let i = 0; i < temps.length; i++) {
        const t = temps[i];
        while (stack.length > 0 && t > stack[stack.length - 1][0]) {
            const [, j] = stack.pop()!;
            result[j] = i - j;
        }
        stack.push([t, i]);
    }
    return result;
}
