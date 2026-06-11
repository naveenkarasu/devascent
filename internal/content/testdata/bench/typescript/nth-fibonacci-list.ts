function fibonacci_list(n: number): number[] {
    if (n === 0) return [];
    if (n === 1) return [0];
    const result: number[] = [0, 1];
    for (let i = 2; i < n; i++) {
        result.push(result[result.length - 1] + result[result.length - 2]);
    }
    return result;
}
