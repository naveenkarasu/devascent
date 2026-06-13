function walls_and_gates(rooms: number[][]): number[][] {
    if (!rooms || rooms.length === 0 || rooms[0].length === 0) return rooms;
    const INF = 2147483647;
    const m = rooms.length;
    const n = rooms[0].length;
    const queue: [number, number][] = [];
    for (let i = 0; i < m; i++) {
        for (let j = 0; j < n; j++) {
            if (rooms[i][j] === 0) {
                queue.push([i, j]);
            }
        }
    }
    let qi = 0;
    while (qi < queue.length) {
        const [i, j] = queue[qi++];
        for (const [di, dj] of [[1,0],[-1,0],[0,1],[0,-1]]) {
            const ni = i + di;
            const nj = j + dj;
            if (ni >= 0 && ni < m && nj >= 0 && nj < n && rooms[ni][nj] === INF) {
                rooms[ni][nj] = rooms[i][j] + 1;
                queue.push([ni, nj]);
            }
        }
    }
    return rooms;
}
