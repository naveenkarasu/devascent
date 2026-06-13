public class Solution {
    public long[] merge_k_lists(long[][] lists) {
        var heap = new SortedSet<(long val, int i, int j)>(
            Comparer<(long val, int i, int j)>.Create((a, b) => {
                int c = a.val.CompareTo(b.val);
                if (c != 0) return c;
                c = a.i.CompareTo(b.i);
                if (c != 0) return c;
                return a.j.CompareTo(b.j);
            })
        );
        if (lists == null) return new long[0];
        for (int i = 0; i < lists.Length; i++)
            if (lists[i] != null && lists[i].Length > 0)
                heap.Add((lists[i][0], i, 0));

        var out_list = new List<long>();
        while (heap.Count > 0) {
            var (val, i, j) = heap.Min;
            heap.Remove((val, i, j));
            out_list.Add(val);
            if (j + 1 < lists[i].Length)
                heap.Add((lists[i][j + 1], i, j + 1));
        }
        return out_list.ToArray();
    }
}
