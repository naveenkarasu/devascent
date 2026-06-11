function count_even_sum_triples(arr: number[], l: number, r: number): number {
    const sub = arr.slice(l - 1, r);
    const e = sub.filter(x => x % 2 === 0).length;
    const o = sub.length - e;
    return Math.floor((e * (e - 1) * (e - 2)) / 6) + Math.floor((o * (o - 1)) / 2) * e;
}
