function smallest_trailing_zero_multiple(a: number, b: number): number {
    function gcd(x: number, y: number): number {
        while (y !== 0) {
            [x, y] = [y, x % y];
        }
        return x;
    }
    const power = Math.pow(10, b);
    return (a / gcd(a, power)) * power;
}
