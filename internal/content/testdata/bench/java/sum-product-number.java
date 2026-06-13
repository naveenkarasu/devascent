class Solution {
    public boolean is_sum_product_number(long n) {
        long digitSum = 0;
        long digitProd = 1;
        long temp = n;
        while (temp > 0) {
            long d = temp % 10;
            digitSum += d;
            digitProd *= d;
            temp /= 10;
        }
        return digitSum * digitProd == n;
    }
}
