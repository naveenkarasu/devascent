function can_complete_circuit(gas, cost) {
    const totalGas = gas.reduce((a, b) => a + b, 0);
    const totalCost = cost.reduce((a, b) => a + b, 0);
    if (totalGas < totalCost) return -1;
    let total = 0, start = 0;
    for (let i = 0; i < gas.length; i++) {
        total += gas[i] - cost[i];
        if (total < 0) {
            start = i + 1;
            total = 0;
        }
    }
    return start;
}
