public class Solution {
    public long can_complete_circuit(long[] gas, long[] cost) {
        long totalGas = 0, totalCost = 0;
        foreach (long g in gas) totalGas += g;
        foreach (long c in cost) totalCost += c;
        if (totalGas < totalCost) return -1;
        long total = 0, start = 0;
        for (int i = 0; i < gas.Length; i++) {
            total += gas[i] - cost[i];
            if (total < 0) { start = i + 1; total = 0; }
        }
        return start;
    }
}
