function max_div_plus_mod(l: number, r: number, a: number): number {
    let best = Math.floor(r / a) + (r % a);
    const multiple = Math.floor(r / a) * a;
    const candidate = multiple - 1;
    if (l <= candidate && candidate <= r) {
        best = Math.max(best, Math.floor(candidate / a) + (candidate % a));
    }
    return best;
}
