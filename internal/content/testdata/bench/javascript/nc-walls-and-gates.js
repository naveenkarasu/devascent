function walls_and_gates(rooms) {
    if (!rooms || !rooms[0]) return rooms;
    const INF = 2147483647;
    const m = rooms.length, n = rooms[0].length;
    const queue = [];
    for (let i = 0; i < m; i++) {
        for (let j = 0; j < n; j++) {
            if (rooms[i][j] === 0) queue.push([i, j]);
        }
    }
    let head = 0;
    while (head < queue.length) {
        const [i, j] = queue[head++];
        for (const [di, dj] of [[1, 0], [-1, 0], [0, 1], [0, -1]]) {
            const ni = i + di, nj = j + dj;
            if (ni >= 0 && ni < m && nj >= 0 && nj < n && rooms[ni][nj] === INF) {
                rooms[ni][nj] = rooms[i][j] + 1;
                queue.push([ni, nj]);
            }
        }
    }
    return rooms;
}
