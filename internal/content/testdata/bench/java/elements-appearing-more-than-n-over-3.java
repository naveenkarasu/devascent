import java.util.*;

class Solution {
    public long[] majority_elements_n3(long[] nums) {
        long cnt1 = 0, cnt2 = 0;
        Long num1 = null, num2 = null;
        for (long n : nums) {
            if (num1 != null && n == num1) {
                cnt1++;
            } else if (num2 != null && n == num2) {
                cnt2++;
            } else if (cnt1 == 0) {
                num1 = n;
                cnt1 = 1;
            } else if (cnt2 == 0) {
                num2 = n;
                cnt2 = 1;
            } else {
                cnt1--;
                cnt2--;
            }
        }
        long threshold = nums.length / 3;
        List<Long> result = new ArrayList<>();
        if (num1 != null) {
            long count = 0;
            for (long x : nums) if (x == num1) count++;
            if (count > threshold) result.add(num1);
        }
        if (num2 != null && !num2.equals(num1)) {
            long count = 0;
            for (long x : nums) if (x == num2) count++;
            if (count > threshold) result.add(num2);
        }
        Collections.sort(result);
        long[] arr = new long[result.size()];
        for (int i = 0; i < result.size(); i++) arr[i] = result.get(i);
        return arr;
    }
}
