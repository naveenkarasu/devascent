using System.Collections.Generic;

public class Solution {
    public long[][] all_permutations(long[] nums) {
        var result = new List<long[]>();
        Permute(nums, result);
        return result.ToArray();
    }

    private void Permute(long[] nums, List<long[]> result) {
        if (nums.Length == 0) {
            result.Add(new long[0]);
            return;
        }
        for (int i = 0; i < nums.Length; i++) {
            long n = nums[i];
            long[] rest = new long[nums.Length - 1];
            int idx = 0;
            for (int j = 0; j < nums.Length; j++) {
                if (j != i) rest[idx++] = nums[j];
            }
            var subPerms = new List<long[]>();
            Permute(rest, subPerms);
            foreach (long[] perm in subPerms) {
                long[] full = new long[perm.Length + 1];
                full[0] = n;
                System.Array.Copy(perm, 0, full, 1, perm.Length);
                result.Add(full);
            }
        }
    }
}
