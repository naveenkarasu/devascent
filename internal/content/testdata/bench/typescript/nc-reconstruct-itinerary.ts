function find_itinerary(tickets: string[][]): string[] {
    const graph: {[key: string]: string[]} = {};
    // Sort tickets in reverse so that pop() yields lexicographically smallest
    tickets.sort((a, b) => {
        if (a[0] !== b[0]) return b[0] < a[0] ? -1 : 1;
        return b[1] < a[1] ? -1 : 1;
    });
    for (const [src, dst] of tickets) {
        if (!(src in graph)) graph[src] = [];
        graph[src].push(dst);
    }

    const result: string[] = [];

    function dfs(airport: string): void {
        while (graph[airport] && graph[airport].length > 0) {
            dfs(graph[airport].pop()!);
        }
        result.push(airport);
    }

    dfs("JFK");
    return result.reverse();
}
