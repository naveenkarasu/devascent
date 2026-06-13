class Solution {
    public long net_salary(long gross) {
        long tax;
        if (gross <= 300000) {
            tax = 0;
        } else if (gross <= 500000) {
            tax = (gross - 300000) * 5 / 100;
        } else if (gross <= 1000000) {
            tax = 200000 * 5 / 100 + (gross - 500000) * 20 / 100;
        } else {
            tax = 200000 * 5 / 100 + 500000 * 20 / 100 + (gross - 1000000) * 30 / 100;
        }
        return gross - tax;
    }
}
