function pacific_atlantic(heights) {
    if (!heights || heights.length === 0) return [];
    const m = heights.length, n = heights[0].length;

    function bfs(starts) {
        const visited = new Set(starts.map(([i, j]) => i * n + j));
        const queue = [...starts];
        let qi = 0;
        while (qi < queue.length) {
            const [i, j] = queue[qi++];
            for (const [di, dj] of [[1,0],[-1,0],[0,1],[0,-1]]) {
                const ni = i + di, nj = j + dj;
                if (ni >= 0 && ni < m && nj >= 0 && nj < n && !visited.has(ni * n + nj) && heights[ni][nj] >= heights[i][j]) {
                    visited.add(ni * n + nj);
                    queue.push([ni, nj]);
                }
            }
        }
        return visited;
    }

    const pacStarts = [];
    for (let j = 0; j < n; j++) pacStarts.push([0, j]);
    for (let i = 1; i < m; i++) pacStarts.push([i, 0]);
    const atlStarts = [];
    for (let j = 0; j < n; j++) atlStarts.push([m-1, j]);
    for (let i = 0; i < m-1; i++) atlStarts.push([i, n-1]);

    const pac = bfs(pacStarts);
    const atl = bfs(atlStarts);

    const result = [];
    for (let i = 0; i < m; i++) {
        for (let j = 0; j < n; j++) {
            if (pac.has(i * n + j) && atl.has(i * n + j)) {
                result.push([i, j]);
            }
        }
    }
    // iterating row-major gives sorted order already
    return result;
}
