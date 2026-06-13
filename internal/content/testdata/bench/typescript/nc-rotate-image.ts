function rotate(matrix: number[][]): number[][] {
    const n = matrix.length;
    // Transpose
    for (let i = 0; i < n; i++) {
        for (let j = i + 1; j < n; j++) {
            const tmp = matrix[i][j];
            matrix[i][j] = matrix[j][i];
            matrix[j][i] = tmp;
        }
    }
    // Reverse each row
    for (const row of matrix) {
        row.reverse();
    }
    return matrix;
}
