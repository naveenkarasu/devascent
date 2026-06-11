function insert_to_minimize(a: string): string {
    if (a[0] === '1') {
        return a[0] + '0' + a.slice(1);
    }
    return '1' + a;
}
