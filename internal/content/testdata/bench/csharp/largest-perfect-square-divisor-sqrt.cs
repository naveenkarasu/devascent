public class Solution {
    public long max_square_divisor_root(long n) {
        long answer = 1;
        long remaining = n;
        long fact = 2;
        while (fact * fact <= n) {
            int count = 0;
            while (remaining % fact == 0) {
                count++;
                remaining /= fact;
                if (count == 2) {
                    answer *= fact;
                    count = 0;
                }
            }
            fact++;
        }
        return answer;
    }
}
