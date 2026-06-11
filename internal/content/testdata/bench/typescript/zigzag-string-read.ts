function zigzag_convert(s: string, num_rows: number): string {
    if (num_rows === 1) return s;
    const rows: string[] = new Array(num_rows).fill('');
    let row = 0, direction = -1;
    for (const ch of s) {
        rows[row] += ch;
        if (row === 0 || row === num_rows - 1) {
            direction = -direction;
        }
        row += direction;
    }
    return rows.join('');
}
