class Solution {
    public long[] sort_three_values(long[] arr, long low_val, long mid_val, long high_val) {
        int lo = 0, mid = 0, hi = arr.length - 1;
        while (mid <= hi) {
            if (arr[mid] == low_val) {
                long tmp = arr[lo]; arr[lo] = arr[mid]; arr[mid] = tmp;
                lo++;
                mid++;
            } else if (arr[mid] == mid_val) {
                mid++;
            } else {
                long tmp = arr[mid]; arr[mid] = arr[hi]; arr[hi] = tmp;
                hi--;
            }
        }
        return arr;
    }
}
