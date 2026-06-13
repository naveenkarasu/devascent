function gcd_lcm(a: number, b: number): number[] {
    function gcd(x: number, y: number): number {
        while (y) {
            const tmp = y;
            y = x % y;
            x = tmp;
        }
        return x;
    }
    const g = gcd(a, b);
    const l = Math.floor(a * b / g);
    return [g, l];
}
