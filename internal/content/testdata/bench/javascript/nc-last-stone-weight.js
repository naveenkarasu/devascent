function last_stone_weight(stones) {
    const heap = [...stones];
    // Use a max-heap simulation: sort descending each iteration
    while (heap.length > 1) {
        heap.sort((a, b) => b - a);
        const y = heap.shift();
        const x = heap.shift();
        if (x !== y) heap.push(y - x);
    }
    return heap.length > 0 ? heap[0] : 0;
}
