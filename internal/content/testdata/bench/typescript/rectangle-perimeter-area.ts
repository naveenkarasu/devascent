function rectangle_perimeter_area(a: number, b: number): number[] {
    if (a <= 0 || b <= 0) return [0, 0];
    return [(a + b) * 2, a * b];
}
