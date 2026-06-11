using System.Collections.Generic;

public class Solution {
    public long[] majority_elements_n3(long[] nums) {
        long cnt1 = 0, cnt2 = 0;
        long? num1 = null, num2 = null;
        foreach (long n in nums) {
            if (num1.HasValue && n == num1.Value) {
                cnt1++;
            } else if (num2.HasValue && n == num2.Value) {
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
        long threshold = nums.Length / 3;
        var result = new List<long>();
        if (num1.HasValue) {
            long count = 0;
            foreach (long x in nums) if (x == num1.Value) count++;
            if (count > threshold) result.Add(num1.Value);
        }
        if (num2.HasValue && num2.Value != num1.GetValueOrDefault(-1)) {
            long count = 0;
            foreach (long x in nums) if (x == num2.Value) count++;
            if (count > threshold) result.Add(num2.Value);
        }
        result.Sort();
        return result.ToArray();
    }
}
