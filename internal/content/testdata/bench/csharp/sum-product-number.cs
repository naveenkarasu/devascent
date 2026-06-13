public class Solution {
    public bool is_sum_product_number(long n) {
        long digitSum = 0, digitProd = 1;
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
