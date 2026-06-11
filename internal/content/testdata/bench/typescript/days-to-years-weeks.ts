function days_breakdown(d: number): number[] {
    const y = Math.floor(d / 365);
    d = d % 365;
    const w = Math.floor(d / 7);
    d = d % 7;
    return [y, w, d];
}
