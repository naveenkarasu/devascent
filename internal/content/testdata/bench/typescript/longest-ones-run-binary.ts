function longest_ones_run(n: number): number {
    const binary = n.toString(2);
    const runs = binary.split('0');
    return Math.max(...runs.map(r => r.length));
}
