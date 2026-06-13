public class Solution {
    public long find_duplicate(long[] nums) {
        long slow = nums[0];
        long fast = nums[(int)nums[0]];
        while (slow != fast) {
            slow = nums[(int)slow];
            fast = nums[(int)nums[(int)fast]];
        }
        long slow2 = 0;
        while (slow != slow2) {
            slow = nums[(int)slow];
            slow2 = nums[(int)slow2];
        }
        return slow;
    }
}
