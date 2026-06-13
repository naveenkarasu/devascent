function rotate_matrix_cw(matrix: number[][]): number[][] {
    const n = matrix.length;
    const result: number[][] = [];
    for (let i = 0; i < n; i++) {
        const row: number[] = [];
        for (let j = 0; j < n; j++) {
            row.push(matrix[n - 1 - j][i]);
        }
        result.push(row);
    }
    return result;
}
