import java.util.*;

class Solution {
    public long[] majority_element_ii(long[] nums) {
        long cnt1 = 0, cnt2 = 0, num1 = 0, num2 = 1;
        for (long n : nums) {
            if (num1 == n) cnt1++;
            else if (num2 == n) cnt2++;
            else if (cnt1 == 0) { num1 = n; cnt1 = 1; }
            else if (cnt2 == 0) { num2 = n; cnt2 = 1; }
            else { cnt1--; cnt2--; }
        }
        cnt1 = 0; cnt2 = 0;
        for (long n : nums) {
            if (n == num1) cnt1++;
            else if (n == num2) cnt2++;
        }
        List<Long> res = new ArrayList<>();
        if (cnt1 > nums.length / 3) res.add(num1);
        if (cnt2 > nums.length / 3) res.add(num2);
        Collections.sort(res);
        long[] out = new long[res.size()];
        for (int i = 0; i < res.size(); i++) out[i] = res.get(i);
        return out;
    }
}
