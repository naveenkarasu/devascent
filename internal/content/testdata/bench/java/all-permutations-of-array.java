import java.util.*;

class Solution {
    public long[][] all_permutations(long[] nums) {
        List<long[]> result = new ArrayList<>();
        permute(nums, result);
        return result.toArray(new long[0][]);
    }

    private void permute(long[] nums, List<long[]> result) {
        if (nums.length == 0) {
            result.add(new long[0]);
            return;
        }
        for (int i = 0; i < nums.length; i++) {
            long n = nums[i];
            long[] rest = new long[nums.length - 1];
            int idx = 0;
            for (int j = 0; j < nums.length; j++) {
                if (j != i) rest[idx++] = nums[j];
            }
            List<long[]> subPerms = new ArrayList<>();
            permute(rest, subPerms);
            for (long[] perm : subPerms) {
                long[] full = new long[perm.length + 1];
                full[0] = n;
                System.arraycopy(perm, 0, full, 1, perm.length);
                result.add(full);
            }
        }
    }
}
