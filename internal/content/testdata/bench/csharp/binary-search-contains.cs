public class Solution {
    public bool contains_element(long[] arr, long key) {
        int lo = 0, hi = arr.Length - 1;
        while (lo <= hi) {
            int mid = (lo + hi) / 2;
            if (arr[mid] == key) return true;
            else if (arr[mid] < key) lo = mid + 1;
            else hi = mid - 1;
        }
        return false;
    }
}
