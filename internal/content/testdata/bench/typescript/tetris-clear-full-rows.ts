function clear_full_rows(board: string[][], empty: string): string[][] {
    const cols = board.length > 0 ? board[0].length : 0;
    const surviving = board.filter(row => row.includes(empty));
    const clearedCount = board.length - surviving.length;
    const blankRows: string[][] = [];
    for (let i = 0; i < clearedCount; i++) {
        blankRows.push(new Array(cols).fill(empty));
    }
    return blankRows.concat(surviving);
}
