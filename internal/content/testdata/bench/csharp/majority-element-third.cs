using System.Linq;

public class Solution {
    public long[] majority_element_ii(long[] nums) {
        long cnt1 = 0, cnt2 = 0, num1 = 0, num2 = 1;
        foreach (long n in nums) {
            if (num1 == n) cnt1++;
            else if (num2 == n) cnt2++;
            else if (cnt1 == 0) { num1 = n; cnt1 = 1; }
            else if (cnt2 == 0) { num2 = n; cnt2 = 1; }
            else { cnt1--; cnt2--; }
        }
        cnt1 = cnt2 = 0;
        foreach (long n in nums) {
            if (n == num1) cnt1++;
            else if (n == num2) cnt2++;
        }
        var res = new List<long>();
        if (cnt1 > nums.Length / 3) res.Add(num1);
        if (cnt2 > nums.Length / 3) res.Add(num2);
        res.Sort();
        return res.ToArray();
    }
}
