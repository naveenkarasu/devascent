function solve(board: string[][]): string[][] {
    if (!board || board.length === 0 || board[0].length === 0) return board;
    const m = board.length;
    const n = board[0].length;
    const safe = new Set<number>();

    function dfs(i: number, j: number): void {
        if (i < 0 || j < 0 || i >= m || j >= n) return;
        if (board[i][j] !== "O" || safe.has(i * n + j)) return;
        safe.add(i * n + j);
        dfs(i+1,j); dfs(i-1,j); dfs(i,j+1); dfs(i,j-1);
    }

    for (let i = 0; i < m; i++) {
        dfs(i, 0); dfs(i, n - 1);
    }
    for (let j = 0; j < n; j++) {
        dfs(0, j); dfs(m - 1, j);
    }

    const result: string[][] = [];
    for (let i = 0; i < m; i++) {
        const row: string[] = [];
        for (let j = 0; j < n; j++) {
            if (board[i][j] === "O" && !safe.has(i * n + j)) {
                row.push("X");
            } else {
                row.push(board[i][j]);
            }
        }
        result.push(row);
    }
    return result;
}
