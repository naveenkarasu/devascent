public class Solution {
    public long[] merge_sorted_arrays(long[] arr1, long[] arr2) {
        var result = new List<long>();
        int i = 0, j = 0;
        while (i < arr1.Length && j < arr2.Length) {
            if (arr1[i] <= arr2[j]) result.Add(arr1[i++]);
            else result.Add(arr2[j++]);
        }
        while (i < arr1.Length) result.Add(arr1[i++]);
        while (j < arr2.Length) result.Add(arr2[j++]);
        return result.ToArray();
    }
}
